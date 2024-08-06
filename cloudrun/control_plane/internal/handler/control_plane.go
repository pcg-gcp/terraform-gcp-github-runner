package handler

import (
	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/config"
	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/gcp"
	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/ghclient"
)

type ControlPlaneHandler struct {
	cfg          *config.Config
	gcpClient    *gcp.Client
	githubClient *ghclient.Client
}

type eventSummaryMessage struct {
	Repository     string `json:"repository"`
	Owner          string `json:"owner"`
	EventType      string `json:"eventType"`
	ID             int64  `json:"id"`
	InstallationID int64  `json:"installationId"`
}

func New(cfg *config.Config, githubClient *ghclient.Client, gcpClient *gcp.Client) (*ControlPlaneHandler, error) {
	return &ControlPlaneHandler{
		cfg:          cfg,
		githubClient: githubClient,
		gcpClient:    gcpClient,
	}, nil
}
