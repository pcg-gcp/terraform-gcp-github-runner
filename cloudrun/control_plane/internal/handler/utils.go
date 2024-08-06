package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
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
	jobStatus, err := h.githubClient.GetJobStatus(m.ID, m.InstallationID, m.Owner, m.Repository, ctx)
	if err != nil {
		err = fmt.Errorf("failed to get jobStatus: %w", err)
		return false, err
	}
	if jobStatus != "queued" {
		slog.Warn("Job no longer queued. No instance will be created")
		return false, nil
	}

	instanceList, err := h.gcpClient.GetInstances()
	if err != nil {
		err = fmt.Errorf("failed to get instance list: %w", err)
		return false, err
	}
	instanceCount := 0
	availableZones, err := h.gcpClient.GetAvailableZones()
	if err != nil {
		err = fmt.Errorf("failed to get available zones: %w", err)
		return false, err
	}
	for _, zone := range availableZones {
		instanceCount += len(instanceList.Items[zone].Instances)
	}
	if instanceCount >= h.cfg.MaxRunnerCount {
		slog.Warn("Already reached max instance count. Scale up not possible", "instance count", instanceCount)
		return false, nil
	}
	return true, nil
}
