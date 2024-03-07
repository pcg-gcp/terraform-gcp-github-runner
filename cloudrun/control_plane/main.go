package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v60/github"
	"github.com/sethvargo/go-envconfig"
	compute "google.golang.org/api/compute/v1"
)

type config struct {
	ProjectID            string `env:"PROJECT_ID,required"`
	Zone                 string `env:"ZONE,required"`
	Region               string `env:"REGION,required"`
	MachineType          string `env:"MACHINE_TYPE,required"`
	ImagePath            string `env:"IMAGE_PATH,required"`
	GithubAppPrivateKey  string `env:"GITHUB_APP_PRIVATE_KEY,required"`
	RunnerServiceAccount string `env:"RUNNER_SERVICE_ACCOUNT,required"`
	Network              string `env:"NETWORK,required"`
	Subnet               string `env:"SUBNET,required"`
	RunnerUser           string `env:"RUNNER_USER,required"`
	RunnerDir            string `env:"RUNNER_DIR,required"`
	StartupScriptURL     string `env:"STARTUP_SCRIPT_URL,required"`
	AppID                int64  `env:"GITHUB_APP_ID,required"`
	Port                 int    `env:"PORT,default=8080"`
	Debug                bool   `env:"DEBUG,default=false"`
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
	http.HandleFunc("/startup", StartRunner)
	addr := fmt.Sprintf(":%d", c.Port)
	slog.Info(fmt.Sprintf("Starting server on %s", addr))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting server: %s", err))
	}
	slog.Info("Server stopped")
}

func StartRunner(w http.ResponseWriter, r *http.Request) {
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

	privateKeyBytes, err := base64.StdEncoding.DecodeString(c.GithubAppPrivateKey)
	if err != nil {
		slog.Error(fmt.Sprintf("Error decoding private key: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	itr, err := ghinstallation.New(http.DefaultTransport, c.AppID, m.InstallationID, privateKeyBytes)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating installation client: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	client := github.NewClient(&http.Client{Transport: itr})

	token, _, err := client.Actions.CreateRegistrationToken(ctx, m.Owner, m.Repository)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating registration token: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Info(fmt.Sprintf("Token: %s", token.GetToken()))

	configItems := []string{
		fmt.Sprintf("--url https://github.com/%s/%s", m.Owner, m.Repository),
		fmt.Sprintf("--token %s", token.GetToken()),
	}
	githubRunnerConfig := strings.Join(configItems, " ")

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
				Network:    "global/networks/" + c.Network,
				Subnetwork: "regions/" + c.Region + "/subnetworks/" + c.Subnet,
			},
		},
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Email: c.RunnerServiceAccount,
				Scopes: []string{
					"https://www.googleapis.com/auth/cloud-platform",
				},
			},
		},
		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{
				{
					Key:   "startup-script-url",
					Value: &c.StartupScriptURL,
				},
				{
					Key:   "github_runner_config",
					Value: &githubRunnerConfig,
				},
				{
					Key:   "runner_user",
					Value: &c.RunnerUser,
				},
				{
					Key:   "runner_dir",
					Value: &c.RunnerDir,
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
