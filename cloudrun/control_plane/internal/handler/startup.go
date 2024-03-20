package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v60/github"
	"google.golang.org/api/compute/v1"
)

func (h *ControlPlaneHandler) StartRunner(w http.ResponseWriter, r *http.Request) {
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
	slog.Info(fmt.Sprintf("Proccesing event for %s/%s", m.Owner, m.Repository))

	ctx := r.Context()

	slog.Debug(fmt.Sprintf("Generating installation client for installation %d", m.InstallationID))

	itr, err := ghinstallation.New(http.DefaultTransport, h.cfg.AppID, m.InstallationID, h.privateKeyBytes)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating installation client: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ghClient := github.NewClient(&http.Client{Transport: itr})

	slog.Debug(fmt.Sprintf("Creating registration token for %s/%s", m.Owner, m.Repository))

	token, _, err := ghClient.Actions.CreateRegistrationToken(ctx, m.Owner, m.Repository)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating registration token: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	configItems := []string{
		fmt.Sprintf("--url https://github.com/%s/%s", m.Owner, m.Repository),
		fmt.Sprintf("--token %s", token.GetToken()),
	}

	if h.cfg.Ephemeral {
		configItems = append(configItems, "--ephemeral")
	}
	githubRunnerConfig := strings.Join(configItems, " ")

	slog.Debug(fmt.Sprintf("Getting instance template %s", h.cfg.InstanceTemplateName))
	template, err := h.computeService.InstanceTemplates.Get(h.cfg.ProjectID, h.cfg.InstanceTemplateName).Do()
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting instance template: %s", err))
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

	slog.Debug(fmt.Sprintf("Creating instance %s", instanceName))

	instance := &compute.Instance{
		Name: instanceName,
		Metadata: &compute.Metadata{
			Items: append(template.Properties.Metadata.Items, &compute.MetadataItems{
				Key:   "github_runner_config",
				Value: &githubRunnerConfig,
			}),
		},
		Labels: map[string]string{
			"ghr-managed": "true",
			"ghr-type":    "repo",
			"ghr-repo":    m.Repository,
			"ghr-owner":   m.Owner,
		},
	}
	createInstanceRequest := h.computeService.Instances.Insert(h.cfg.ProjectID, h.cfg.Zone, instance)
	createInstanceRequest = createInstanceRequest.SourceInstanceTemplate(template.SelfLink)
	op, err := createInstanceRequest.Do()
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating instance: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	for op.Status != "DONE" {
		time.Sleep(1 * time.Second)
		op, err = h.computeService.ZoneOperations.Get(h.cfg.ProjectID, h.cfg.Zone, op.Name).Do()
		if err != nil {
			slog.Error(fmt.Sprintf("Error getting operation status: %s", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	slog.Info(fmt.Sprintf("Instance %s created", instanceName))

	fmt.Fprint(w, "Success!")
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
