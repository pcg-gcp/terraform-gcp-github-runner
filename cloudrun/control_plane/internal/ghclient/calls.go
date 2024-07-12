package ghclient

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/go-github/v63/github"
	"google.golang.org/api/compute/v1"
)

func GenerateRunnerConfig(installationID int64, owner, repository, instanceName string, ctx context.Context) (*compute.MetadataItems, string, error) {
	client, err := getClient(installationID)
	if err != nil {
		return nil, "", err
	}

	var configMetadata *compute.MetadataItems
	var useJitConfigStr string
	if cfg.Ephemeral && cfg.UseJitConfig {
		useJitConfigStr = "true"
		workfolder := "_work"
		jitConfig, _, err := client.Actions.GenerateRepoJITConfig(ctx, owner, repository, &github.GenerateJITConfigRequest{
			Labels:        []string{"self-hosted", "ephemeral"},
			Name:          instanceName,
			WorkFolder:    &workfolder,
			RunnerGroupID: 1,
		})
		if err != nil {
			err = fmt.Errorf("failed to generate jitConfig: %w", err)
			return nil, "", err
		}
		encodedJITConfig := jitConfig.GetEncodedJITConfig()
		configMetadata = &compute.MetadataItems{
			Key:   "github_runner_config",
			Value: &encodedJITConfig,
		}
	} else {
		useJitConfigStr = "false"
		token, _, err := client.Actions.CreateRegistrationToken(ctx, owner, repository)
		if err != nil {
			err = fmt.Errorf("failed to generate registration token: %w", err)
			return nil, "", err
		}

		configItems := []string{
			fmt.Sprintf("--url https://github.com/%s/%s", owner, repository),
			fmt.Sprintf("--token %s", token.GetToken()),
		}

		if cfg.Ephemeral {
			configItems = append(configItems, "--ephemeral")
		}
		githubRunnerConfig := strings.Join(configItems, " ")
		configMetadata = &compute.MetadataItems{
			Key:   "github_runner_config",
			Value: &githubRunnerConfig,
		}
	}
	return configMetadata, useJitConfigStr, nil
}

func RemoveRunnerForInstance(instance *compute.Instance, ctx context.Context) (bool, error) {
	var installationId int64

	repo := instance.Labels["ghr-repo"]
	owner := instance.Labels["ghr-owner"]
	runnerType := instance.Labels["ghr-type"]

	appsClient := getAppsClient()

	switch runnerType {
	case "repo":
		installation, _, err := appsClient.Apps.FindRepositoryInstallation(ctx, owner, repo)
		if err != nil {
			err = fmt.Errorf("failed to find installation: %w", err)
			return false, err
		}
		installationId = installation.GetID()
	case "org":
		installation, _, err := appsClient.Apps.FindOrganizationInstallation(ctx, owner)
		if err != nil {
			err = fmt.Errorf("failed to find installation: %w", err)
			return false, err
		}
		installationId = installation.GetID()
	}
	client, err := getClient(installationId)
	if err != nil {
		err = fmt.Errorf("failed to create client: %w", err)
		return false, err

	}

	var runners *github.Runners
	switch runnerType {
	case "repo":
		runners, _, err = client.Actions.ListRunners(ctx, owner, repo, nil)
		if err != nil {
			err = fmt.Errorf("failed to list runners: %w", err)
			return false, err
		}
	case "org":
		runners, _, err = client.Actions.ListOrganizationRunners(ctx, owner, nil)
		if err != nil {
			err = fmt.Errorf("failed to list runners: %w", err)
			return false, err
		}
	}
	var runner *github.Runner
	for _, itRunner := range runners.Runners {
		if strings.HasSuffix(itRunner.GetName(), instance.Name) {
			runner = itRunner
			break
		}
	}
	if runner == nil {
		slog.Info(fmt.Sprintf("Runner %s not found", instance.Name))
		return true, nil
	}
	if runner.GetBusy() {
		slog.Info(fmt.Sprintf("Runner %s is busy. Skipping", instance.Name))
		return false, nil
	}
	switch runnerType {
	case "repo":
		_, err = client.Actions.RemoveRunner(ctx, owner, repo, runner.GetID())
		if err != nil {
			err = fmt.Errorf("failed to remove runners: %w", err)
			return false, err
		}
	case "org":
		_, err = client.Actions.RemoveOrganizationRunner(ctx, owner, runner.GetID())
		if err != nil {
			err = fmt.Errorf("failed to remove runners: %w", err)
			return false, err
		}
	}
	return true, nil
}

func GetJobStatus(jobID, installationID int64, owner, repo string, ctx context.Context) (string, error) {
	client, err := getClient(installationID)
	if err != nil {
		err = fmt.Errorf("failed to create client: %w", err)
		return "", err
	}

	job, _, err := client.Actions.GetWorkflowJobByID(ctx, owner, repo, jobID)
	if err != nil {
		err = fmt.Errorf("failed to get workflow job: %w", err)
		return "", err
	}
	return job.GetStatus(), nil
}
