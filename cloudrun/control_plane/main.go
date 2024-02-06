package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/sethvargo/go-envconfig"
)

type ControlPlane struct{}

type config struct {
	Port  int  `env:"PORT,default=8080"`
	Debug bool `env:"DEBUG,default=false"`
}

var c config

func main() {
	c = config{}
	if err := envconfig.Process(context.Background(), &c); err != nil {
		panic(err)
	}
	logLevel := slog.LevelInfo
	if c.Debug {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.MessageKey:
				a.Key = "message"
			case slog.SourceKey:
				a.Key = "logging.googleapis.com/sourceLocation"
			case slog.LevelKey:
				a.Key = "severity"
			}
			return a
		},
	})))
	h := &ControlPlane{}
	http.Handle("/", h)
	addr := fmt.Sprintf(":%d", c.Port)
	slog.Info(fmt.Sprintf("Starting server on %s", addr))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting server: %s", err))
	}
	slog.Info("Server stopped")
}

func (*ControlPlane) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info(fmt.Sprintf("Received request: %v", r))
	fmt.Fprint(w, "Success!")
}
