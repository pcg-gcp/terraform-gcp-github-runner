package gcp

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/config"
	"google.golang.org/api/compute/v1"
)

type Client struct {
	computeService      *compute.Service
	secretManagerClient *secretmanager.Client
	cfg                 *config.Config
}

func NewClient(cfg *config.Config) (*Client, error) {
	computeService, err := compute.NewService(context.Background())
	if err != nil {
		return nil, err
	}
	secretManagerClient, err := secretmanager.NewClient(context.Background())
	if err != nil {
		return nil, err
	}
	return &Client{
		computeService:      computeService,
		secretManagerClient: secretManagerClient,
		cfg:                 cfg,
	}, nil
}
