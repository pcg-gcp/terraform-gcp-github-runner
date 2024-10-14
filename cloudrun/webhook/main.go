package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/google/go-github/v66/github"
	"github.com/sethvargo/go-envconfig"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GitHubEventMonitor struct {
	cfg              *config
	webhookSecretKey []byte
}

type config struct {
	WebhookSecretKey      string   `env:"WEBHOOK_SECRET_KEY,required"`
	TaskQueuePath         string   `env:"TASK_QUEUE_PATH,required"`
	InvokerServiceAccount string   `env:"INVOKER_SERVICE_ACCOUNT,required"`
	ControlPlaneURL       string   `env:"CONTROL_PLANE_URL,required"`
	RunnerLabels          []string `env:"RUNNER_LABELS,required"`
	DelaySeconds          int      `env:"DELAY_SECONDS,required"`
	Port                  int      `env:"PORT,default=8080"`
	Debug                 bool     `env:"DEBUG,default=false"`
}

type eventSummaryMessage struct {
	Repository     string `json:"repository"`
	Owner          string `json:"owner"`
	EventType      string `json:"eventType"`
	ID             int64  `json:"id"`
	InstallationID int64  `json:"installationId"`
}

func main() {
	c := config{}
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
		cfg:              &c,
	}
	http.Handle("POST /webhook", &s)
	addr := fmt.Sprintf(":%d", c.Port)
	slog.Info(fmt.Sprintf("Starting server on %s", addr))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting server: %s", err))
		os.Exit(1)
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
		action := e.GetAction()
		slog.Info(fmt.Sprintf("Processing Github event %d for %s", e.WorkflowJob.GetID(), e.GetRepo().GetFullName()))
		if action != "queued" {
			slog.Info(fmt.Sprintf("Ignoring event with the action '%s'", action))
			w.WriteHeader(http.StatusAccepted)
			return
		}

		if !hasRequiredLabels(e.WorkflowJob.Labels, s.cfg.RunnerLabels) {
			slog.Info(fmt.Sprintf("Ignoring event for %s because it does not have the required labels", e.GetRepo().GetFullName()))
			w.WriteHeader(http.StatusAccepted)
			return
		}

		ctx := r.Context()
		client, err := cloudtasks.NewClient(ctx)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to create client: %v", err))
			http.Error(w, "Failed to queue event", http.StatusInternalServerError)
			return
		}
		defer client.Close()

		scheduleTime := time.Now().Add(time.Second * time.Duration(s.cfg.DelaySeconds))

		messageStruct := eventSummaryMessage{
			Repository:     e.GetRepo().GetName(),
			Owner:          e.GetRepo().GetOwner().GetLogin(),
			EventType:      eventType,
			ID:             e.WorkflowJob.GetID(),
			InstallationID: e.GetInstallation().GetID(),
		}

		message, err := json.Marshal(messageStruct)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to marshal message: %v", err))
			http.Error(w, "Failed to queue event", http.StatusInternalServerError)
			return
		}

		req := &taskspb.CreateTaskRequest{
			Parent: s.cfg.TaskQueuePath,
			Task: &taskspb.Task{
				ScheduleTime: timestamppb.New(scheduleTime),
				MessageType: &taskspb.Task_HttpRequest{
					HttpRequest: &taskspb.HttpRequest{
						HttpMethod: taskspb.HttpMethod_POST,
						Url:        s.cfg.ControlPlaneURL,
						Headers: map[string]string{
							"Content-Type": "application/json",
						},
						AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
							OidcToken: &taskspb.OidcToken{
								ServiceAccountEmail: s.cfg.InvokerServiceAccount,
							},
						},
						Body: message,
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
		slog.Info(fmt.Sprintf("Queued event %d for %s", e.WorkflowJob.GetID(), e.GetRepo().GetFullName()))
	default:
		slog.Info(fmt.Sprintf("Ignoring event type %s", eventType))
	}
}

func hasRequiredLabels(requiredLabels []string, runnerLabels []string) bool {
	for _, requiredLabel := range requiredLabels {
		hasLabel := false
		for _, runnerLabel := range runnerLabels {
			if requiredLabel == runnerLabel {
				hasLabel = true
				break
			}
		}
		if !hasLabel {
			return false
		}
	}
	return true
}
