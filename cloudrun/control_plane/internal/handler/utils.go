package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/ghclient"
)

func parseStartUpRequest(request *http.Request) (*eventSummaryMessage, error) {
	if request.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("invalid content type %s", request.Header.Get("Content-Type"))
	}
	var message eventSummaryMessage
	if err := json.NewDecoder(request.Body).Decode(&message); err != nil {
		return nil, fmt.Errorf("error decoding request body: %s", err)
	}

	return &message, nil
}

func (h *ControlPlaneHandler) makeStartUpDecision(m *eventSummaryMessage, ctx context.Context) (bool, error) {
	jobStatus, err := ghclient.GetJobStatus(m.ID, m.InstallationID, m.Owner, m.Repository, ctx)
	if err != nil {
		err = fmt.Errorf("failed to get jobStatus: %w", err)
		return false, err
	}
	if jobStatus != "queued" {
		slog.Warn("Job no longer queued. No instance will be created")
		return false, nil
	}

	instanceList, err := h.computeService.Instances.AggregatedList(h.cfg.ProjectID).Filter("labels.ghr-managed=true").Do()
	if err != nil {
		err = fmt.Errorf("failed to get instance list: %w", err)
		return false, err
	}
	if len(instanceList.Items) >= h.cfg.MaxRunnerCount {
		slog.Warn("Already reached max instance count. Scale up not possible", "instance count", len(instanceList.Items))
		return false, nil
	}
	return true, nil
}
