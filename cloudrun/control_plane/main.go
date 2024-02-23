package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
	ProjectID   string `env:"PROJECT_ID,required"`
	Zone        string `env:"ZONE,required"`
	MachineType string `env:"MACHINE_TYPE,required"`
	ImagePath   string `env:"IMAGE_PATH,required"`
	Port        int    `env:"PORT,default=8080"`
	Debug       bool   `env:"DEBUG,default=false"`
}

var c config

type eventSummaryMessage struct {
	Repository     string `json:"repository"`
	Owner          string `json:"owner"`
	EventType      string `json:"eventType"`
	ID             int64  `json:"id"`
	InstallationID int64  `json:"installationId"`
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
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
		slog.Error(fmt.Sprintf("Invalid request method: %s", r.Method))
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error(fmt.Sprintf("Invalid request content type: %s", r.Header.Get("Content-Type")))
		http.Error(w, "Invalid request content type", http.StatusBadRequest)
		return
	}
	var m eventSummaryMessage
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		slog.Error(fmt.Sprintf("Error decoding request body: %s", err))
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

	id, err := randomHex(8)
	if err != nil {
		slog.Error(fmt.Sprintf("Error generating instance ID: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	instanceName := "ghr-" + id

	instance := &compute.Instance{
		Name:        instanceName,
		MachineType: "/zones/" + c.Zone + "/machineTypes/" + c.MachineType,
		Zone:        c.Zone,
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: c.ImagePath,
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
	op, err := service.Instances.Insert(c.ProjectID, c.Zone, instance).Do()
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating instance: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	for op.Status != "DONE" {
		time.Sleep(1 * time.Second)
		op, err = service.ZoneOperations.Get(c.ProjectID, c.Zone, op.Name).Do()
		if err != nil {
			slog.Error(fmt.Sprintf("Error getting operation status: %s", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	fmt.Fprint(w, "Success!")
}
