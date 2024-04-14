package ghclient

import (
	"encoding/base64"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v61/github"
	"github.com/pcg-gcp/terraform-gcp-github-runner/cloudrun/control_plane/internal/config"
)

var (
	privateKeyBytes []byte
	cfg             *config.Config
	appsClient      *github.Client
)

func CreateClient(installationID int64) (*github.Client, error) {
	itr, err := ghinstallation.New(http.DefaultTransport, cfg.AppID, installationID, privateKeyBytes)
	if err != nil {
		return nil, err
	}
	return github.NewClient(&http.Client{Transport: itr}), nil
}

func GetAppsClient() *github.Client {
	return appsClient
}

func Init(cfg *config.Config) error {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(cfg.GithubAppPrivateKey)
	if err != nil {
		return err
	}
	itr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, cfg.AppID, privateKeyBytes)
	if err != nil {
		return err
	}
	appsClient = github.NewClient(&http.Client{Transport: itr})
	return nil
}
