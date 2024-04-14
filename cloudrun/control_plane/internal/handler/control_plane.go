package handler

import (
	"context"

	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/config"
	"google.golang.org/api/compute/v1"
)

type ControlPlaneHandler struct {
	cfg            *config.Config
	computeService *compute.Service
}

type eventSummaryMessage struct {
	Repository     string `json:"repository"`
	Owner          string `json:"owner"`
	EventType      string `json:"eventType"`
	ID             int64  `json:"id"`
	InstallationID int64  `json:"installationId"`
}

func New(cfg *config.Config) (*ControlPlaneHandler, error) {
	computeService, err := compute.NewService(context.Background())
	if err != nil {
		return nil, err
	}
	handler := &ControlPlaneHandler{cfg: cfg, computeService: computeService}
	return handler, nil
}
