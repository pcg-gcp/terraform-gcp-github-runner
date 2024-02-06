package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/google/go-github/v58/github"
	"github.com/sethvargo/go-envconfig"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GitHubEventMonitor struct {
	webhookSecretKey []byte
}

type config struct {
	WebhookSecretKey      string `env:"WEBHOOK_SECRET_KEY,required"`
	ProjectID             string `env:"PROJECT_ID,required"`
	TaskQueuePath         string `env:"TASK_QUEUE_PATH,required"`
	InvokerServiceAccount string `env:"INVOKER_SERVICE_ACCOUNT,required"`
	ControlPlaneURL       string `env:"CONTROL_PLANE_URL,required"`
	Port                  int    `env:"PORT,default=8080"`
	Debug                 bool   `env:"DEBUG,default=false"`
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
	s := GitHubEventMonitor{
		webhookSecretKey: []byte(c.WebhookSecretKey),
	}
	http.Handle("/webhook", &s)
	addr := fmt.Sprintf(":%d", c.Port)
	slog.Info(fmt.Sprintf("Starting server on %s", addr))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting server: %s", err))
	}
	slog.Info("Server stopped")
}

func (s *GitHubEventMonitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Debug(fmt.Sprintf("Received request from %s", r.RemoteAddr))
	payload, err := github.ValidatePayload(r, s.webhookSecretKey)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to validate webhook: %v", err))
		http.Error(w, "Invalid Key", http.StatusBadRequest)
		return
	}
	eventType := github.WebHookType(r)
	event, err := github.ParseWebHook(eventType, payload)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to parse webhook: %v", err))
		http.Error(w, "Invalid Payload", http.StatusBadRequest)
		return
	}
	switch e := event.(type) {
	case *github.WorkflowJobEvent:
		slog.Info(fmt.Sprintf("Processing Github event '%s' for %s", e.GetAction(), e.GetRepo().GetFullName()))
		ctx := r.Context()
		client, err := cloudtasks.NewClient(ctx)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to create client: %v", err))
			http.Error(w, "Failed to queue event", http.StatusInternalServerError)
			return
		}
		defer client.Close()

		scheduleTime := time.Now().Add(time.Second * 10)

		req := &taskspb.CreateTaskRequest{
			Parent: c.TaskQueuePath,
			Task: &taskspb.Task{
				ScheduleTime: timestamppb.New(scheduleTime),
				MessageType: &taskspb.Task_HttpRequest{
					HttpRequest: &taskspb.HttpRequest{
						HttpMethod: taskspb.HttpMethod_POST,
						Url:        c.ControlPlaneURL,
						Headers: map[string]string{
							"Content-Type": r.Header.Get("Content-Type"),
						},
						AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
							OidcToken: &taskspb.OidcToken{
								ServiceAccountEmail: c.InvokerServiceAccount,
							},
						},
						Body: payload,
					},
				},
			},
		}
		_, err = client.CreateTask(ctx, req)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to create task: %v", err))
			http.Error(w, "Failed to queue event", http.StatusInternalServerError)
			return
		}
	default:
		slog.Info(fmt.Sprintf("Ignoring event type %s", eventType))
	}
}
