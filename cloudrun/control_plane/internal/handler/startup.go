package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/gcp"
	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/ghclient"
	"google.golang.org/api/compute/v1"
)

func (h *ControlPlaneHandler) StartRunner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m, err := parseStartUpRequest(r)
	if err != nil {
		slog.Error(fmt.Sprintf("Error validating request: %s", err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	slog.Info(fmt.Sprintf("Proccesing event for %s/%s", m.Owner, m.Repository))

	if ok, err := h.makeStartUpDecision(m, ctx); !ok {
		slog.Info(fmt.Sprintf("Ignoring event for %s/%s", m.Owner, m.Repository))
		fmt.Fprint(w, "Ignored")
		return
	} else if err != nil {
		slog.Error(fmt.Sprintf("Error making startup decision: %s", err))
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

	slog.Debug(fmt.Sprintf("Generating installation client for installation %d", m.InstallationID))

	if h.cfg.UseOrgRunners {
		slog.Debug(fmt.Sprintf("Creating repository registration token for %s/%s", m.Owner, m.Repository))
	} else {
		slog.Debug(fmt.Sprintf("Creating organization registration token for %s", m.Owner))
	}

	githubRunnerConfig, useJitConfigStr, err := ghclient.GenerateRunnerConfig(m.InstallationID, m.Owner, m.Repository, instanceName, h.cfg.UseOrgRunners, ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("Error generating Runner Config: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Debug(fmt.Sprintf("Getting instance template %s", h.cfg.InstanceTemplateName))
	template, err := h.computeService.InstanceTemplates.Get(h.cfg.ProjectID, h.cfg.InstanceTemplateName).Do()
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting instance template: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Debug(fmt.Sprintf("Creating config secret for instance %s", instanceName))
	secretName := fmt.Sprintf("%s-config", instanceName)
	if len(template.Properties.ServiceAccounts) == 0 {
		slog.Error(fmt.Sprintf("No service account found in instance template: %s", h.cfg.InstanceTemplateName))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	serviceAccount := template.Properties.ServiceAccounts[0].Email
	err = gcp.AddSecret(h.cfg.ProjectID, secretName, githubRunnerConfig, serviceAccount, ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating secret: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Debug(fmt.Sprintf("Creating instance %s", instanceName))

	instance := &compute.Instance{
		Name: instanceName,
		Metadata: &compute.Metadata{
			Items: append(template.Properties.Metadata.Items, &compute.MetadataItems{
				Key:   "use_jit_config",
				Value: &useJitConfigStr,
			}, &compute.MetadataItems{
				Key:   "config_secret",
				Value: &secretName,
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
