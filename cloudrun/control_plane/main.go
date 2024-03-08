package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/config"
	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/handler"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	cfg := config.Config{}
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		panic(err)
	}
	logLevel := slog.LevelInfo
	if cfg.Debug {
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
	handler, err := handler.New(&cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating handler: %s", err))
		os.Exit(1)
	}
	http.HandleFunc("/startup", handler.StartRunner)
	http.HandleFunc("/shutdown", handler.StopRunner)
	addr := fmt.Sprintf(":%d", cfg.Port)
	slog.Info(fmt.Sprintf("Starting server on %s", addr))
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting server: %s", err))
		os.Exit(1)
	}
	slog.Info("Server stopped")
}
