package ghclient

import (
	"encoding/base64"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v62/github"
	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/config"
)

var (
	privateKeyBytes []byte
	cfg             *config.Config
	appsClient      *github.Client
	clients         map[int64]*github.Client
)

func GetClient(installationID int64) (*github.Client, error) {
	if client, ok := clients[installationID]; ok {
		return client, nil
	}
	itr, err := ghinstallation.New(http.DefaultTransport, cfg.AppID, installationID, privateKeyBytes)
	if err != nil {
		return nil, err
	}
	client := github.NewClient(&http.Client{Transport: itr})
	clients[installationID] = client
	return client, nil
}

func GetAppsClient() *github.Client {
	return appsClient
}

func Init(config *config.Config) error {
	cfg = config
	var err error
	privateKeyBytes, err = base64.StdEncoding.DecodeString(cfg.GithubAppPrivateKey)
	if err != nil {
		return err
	}
	itr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, cfg.AppID, privateKeyBytes)
	if err != nil {
		return err
	}
	appsClient = github.NewClient(&http.Client{Transport: itr})

	clients = make(map[int64]*github.Client)
	return nil
}
