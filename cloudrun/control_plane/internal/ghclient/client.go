package ghclient

import (
	"encoding/base64"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v79/github"
	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/config"
)

type Client struct {
	cfg             *config.Config
	appsClient      *github.Client
	clients         map[int64]*github.Client
	privateKeyBytes []byte
}

func NewClient(cfg *config.Config) (*Client, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(cfg.GithubAppPrivateKey)
	if err != nil {
		return nil, err
	}
	itr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, cfg.AppID, privateKeyBytes)
	if err != nil {
		return nil, err
	}
	appsClient := github.NewClient(&http.Client{Transport: itr})

	clients := make(map[int64]*github.Client)
	return &Client{
		cfg:             cfg,
		appsClient:      appsClient,
		clients:         clients,
		privateKeyBytes: privateKeyBytes,
	}, nil
}

func (c *Client) getClient(installationID int64) (*github.Client, error) {
	if client, ok := c.clients[installationID]; ok {
		return client, nil
	}
	itr, err := ghinstallation.New(http.DefaultTransport, c.cfg.AppID, installationID, c.privateKeyBytes)
	if err != nil {
		return nil, err
	}
	client := github.NewClient(&http.Client{Transport: itr})
	c.clients[installationID] = client
	return client, nil
}
