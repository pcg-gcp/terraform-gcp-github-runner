package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/v61/github"
	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/ghclient"
	"google.golang.org/api/compute/v1"
)

func (h *ControlPlaneHandler) StartRunner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m, err := ValidateStartUpRequest(r)
	if err != nil {
		slog.Error(fmt.Sprintf("Error validating request: %s", err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	slog.Info(fmt.Sprintf("Proccesing event for %s/%s", m.Owner, m.Repository))

	ghClient, err := ghclient.CreateClient(m.InstallationID)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating client: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if ok, err := MakeStartUpDecision(m); !ok {
		slog.Info(fmt.Sprintf("Ignoring event for %s/%s", m.Owner, m.Repository))
		fmt.Fprint(w, "Ignored")
		w.WriteHeader(http.StatusOK)
		return
	} else if err != nil {
		slog.Error(fmt.Sprintf("Error making decision: %s", err))
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

	slog.Debug(fmt.Sprintf("Creating registration token for %s/%s", m.Owner, m.Repository))

	var configMetadata *compute.MetadataItems
	var useJitConfigStr string
	if h.cfg.Ephemeral && h.cfg.UseJitConfig {
		useJitConfigStr = "true"
		workfolder := "_work"
		jitConfig, _, err := ghClient.Actions.GenerateRepoJITConfig(ctx, m.Owner, m.Repository, &github.GenerateJITConfigRequest{
			Labels:        []string{"self-hosted", "ephemeral"},
			Name:          instanceName,
			WorkFolder:    &workfolder,
			RunnerGroupID: 1,
		})
		if err != nil {
			slog.Error(fmt.Sprintf("Error creating JIT config: %s", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		encodedJITConfig := jitConfig.GetEncodedJITConfig()
		configMetadata = &compute.MetadataItems{
			Key:   "github_runner_config",
			Value: &encodedJITConfig,
		}
	} else {
		useJitConfigStr = "false"
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
		configMetadata = &compute.MetadataItems{
			Key:   "github_runner_config",
			Value: &githubRunnerConfig,
		}
	}

	slog.Debug(fmt.Sprintf("Getting instance template %s", h.cfg.InstanceTemplateName))
	template, err := h.computeService.InstanceTemplates.Get(h.cfg.ProjectID, h.cfg.InstanceTemplateName).Do()
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting instance template: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Debug(fmt.Sprintf("Creating instance %s", instanceName))

	instance := &compute.Instance{
		Name: instanceName,
		Metadata: &compute.Metadata{
			Items: append(template.Properties.Metadata.Items, configMetadata, &compute.MetadataItems{
				Key:   "use_jit_config",
				Value: &useJitConfigStr,
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
