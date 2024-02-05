package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"slices"

	"cloud.google.com/go/pubsub"
	"github.com/google/go-github/v58/github"
	"github.com/sethvargo/go-envconfig"
)

type GitHubEventMonitor struct {
	webhookSecretKey []byte
}

type config struct {
	WebhookSecretKey string `env:"WEBHOOK_SECRET_KEY,required"`
	PubSubProjectID  string `env:"PUBSUB_PROJECT_ID,required"`
	Port             int    `env:"PORT,default=8080"`
	Debug            bool   `env:"DEBUG,default=false"`
}

var allowedEventTypes = []string{
	"workflow_job",
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
		Level: logLevel,
	})))
	s := GitHubEventMonitor{
		webhookSecretKey: []byte(c.WebhookSecretKey),
	}
	http.Handle("/webhook", &s)
	addr := fmt.Sprintf(":%d", c.Port)
	http.ListenAndServe(addr, nil)
}

func (s *GitHubEventMonitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Debug(fmt.Sprintf("Received request from %s", r.RemoteAddr))
	payload, err := github.ValidatePayload(r, s.webhookSecretKey)
	if err != nil {
		http.Error(w, "Invalid Key", http.StatusBadRequest)
		return
	}
	eventType := github.WebHookType(r)
	if !slices.Contains(allowedEventTypes, eventType) {
		http.Error(w, "Invalid Event Type", http.StatusBadRequest)
		return
	}
	event, err := github.ParseWebHook(eventType, payload)
	if err != nil {
		http.Error(w, "Invalid Payload", http.StatusBadRequest)
		return
	}
	switch e := event.(type) {
	case *github.WorkflowJobEvent:
		slog.Info(fmt.Sprintf("Processing Github event %s for %s", e.GetAction(), e.GetRepo().GetFullName()))
	}
	client, err := pubsub.NewClient(context.Background(), c.PubSubProjectID)
	if err != nil {
		http.Error(w, "Error creating pubsub client", http.StatusInternalServerError)
		return
	}

	topic := client.Topic("github-events")

	result := topic.Publish(context.Background(), &pubsub.Message{
		Data: payload,
	})
	_, err = result.Get(context.Background())
	if err != nil {
		http.Error(w, "Error publishing message", http.StatusInternalServerError)
		return
	}
}
