package handler

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v60/github"
	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/config"
	"google.golang.org/api/compute/v1"
)

type ControlPlaneHandler struct {
	cfg             *config.Config
	appsClient      *github.Client
	computeService  *compute.Service
	privateKeyBytes []byte
}

type eventSummaryMessage struct {
	Repository     string `json:"repository"`
	Owner          string `json:"owner"`
	EventType      string `json:"eventType"`
	ID             int64  `json:"id"`
	InstallationID int64  `json:"installationId"`
}

func New(cfg *config.Config) (*ControlPlaneHandler, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(cfg.GithubAppPrivateKey)
	if err != nil {
		return nil, err
	}
	appsItr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, cfg.AppID, privateKeyBytes)
	if err != nil {
		return nil, err
	}
	appsClient := github.NewClient(&http.Client{Transport: appsItr})

	computeService, err := compute.NewService(context.Background())
	if err != nil {
		return nil, err
	}
	handler := &ControlPlaneHandler{cfg: cfg, appsClient: appsClient, computeService: computeService, privateKeyBytes: privateKeyBytes}
	return handler, nil
}
