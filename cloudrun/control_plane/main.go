package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sethvargo/go-envconfig"
	compute "google.golang.org/api/compute/v1"
)

type ControlPlane struct{}

type config struct {
	Port  int  `env:"PORT,default=8080"`
	Debug bool `env:"DEBUG,default=false"`
}

var c config

type eventSummaryMessage struct {
	Repository     string `json:"repository"`
	Owner          string `json:"owner"`
	EventType      string `json:"eventType"`
	ID             int64  `json:"id"`
	InstallationID int64  `json:"installationId"`
}

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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid request content type", http.StatusBadRequest)
		return
	}
	var m eventSummaryMessage
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	service, err := compute.NewService(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating compute service: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	project := "cw-td-sandbox"

	instance := &compute.Instance{
		Name:        "test-instance",
		MachineType: "/zones/europe-west1-b/machineTypes/e2-micro",
		Zone:        "europe-west1-b",
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-12",
				},
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Email: "default",
				Scopes: []string{
					"https://www.googleapis.com/auth/cloud-platform",
				},
			},
		},
	}
	op, err := service.Instances.Insert(project, "europe-west1-b", instance).Do()
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating instance: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	for op.Status != "DONE" {
		time.Sleep(1 * time.Second)
		op, err = service.ZoneOperations.Get(project, "europe-west1-b", op.Name).Do()
		if err != nil {
			slog.Error(fmt.Sprintf("Error getting operation status: %s", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	fmt.Fprint(w, "Success!")
}
