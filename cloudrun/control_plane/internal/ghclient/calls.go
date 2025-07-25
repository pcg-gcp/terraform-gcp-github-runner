package ghclient

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/go-github/v73/github"
	"google.golang.org/api/compute/v1"
)

func (c *Client) generateJitConfig(owner, repository, instanceName string, useOrgRunners bool, ctx context.Context, client *github.Client) (string, error) {
	var jitConfig *github.JITRunnerConfig
	workfolder := "_work"
	var err error

	if !useOrgRunners {
		jitConfig, _, err = client.Actions.GenerateRepoJITConfig(ctx, owner, repository, &github.GenerateJITConfigRequest{
			Labels:        c.cfg.RunnerLabels,
			Name:          instanceName,
			WorkFolder:    &workfolder,
			RunnerGroupID: 1,
		})
		if err != nil {
			err = fmt.Errorf("failed to generate jitConfig: %w", err)
			return "", err
		}
	} else {
		jitConfig, _, err = client.Actions.GenerateOrgJITConfig(ctx, owner, &github.GenerateJITConfigRequest{
			Labels:        c.cfg.RunnerLabels,
			Name:          instanceName,
			WorkFolder:    &workfolder,
			RunnerGroupID: 1,
		})
		if err != nil {
			err = fmt.Errorf("failed to generate jitConfig: %w", err)
			return "", err
		}
	}

	encodedJITConfig := jitConfig.GetEncodedJITConfig()
	return encodedJITConfig, nil
}

func (c *Client) generateStandardConfig(owner, repository string, useOrgRunners bool, ctx context.Context, client *github.Client) (string, error) {
	configItems := []string{fmt.Sprintf("--labels %s", strings.Join(c.cfg.RunnerLabels, ","))}
	var token *github.RegistrationToken
	var err error

	if !useOrgRunners {
		token, _, err = client.Actions.CreateRegistrationToken(ctx, owner, repository)
		if err != nil {
			err = fmt.Errorf("failed to generate registration token: %w", err)
			return "", err
		}

		configItems = append(configItems, fmt.Sprintf("--url https://github.com/%s/%s", owner, repository))
	} else {
		token, _, err = client.Actions.CreateOrganizationRegistrationToken(ctx, owner)
		if err != nil {
			err = fmt.Errorf("failed to generate registration token: %w", err)
			return "", err
		}
		configItems = append(configItems, fmt.Sprintf("--url https://github.com/%s", owner))
	}

	configItems = append(configItems, fmt.Sprintf("--token %s", token.GetToken()))

	if c.cfg.Ephemeral {
		configItems = append(configItems, "--ephemeral")
	}
	githubRunnerConfig := strings.Join(configItems, " ")
	return githubRunnerConfig, nil
}

func (c *Client) GenerateRunnerConfig(installationID int64, owner, repository, instanceName string, useOrgRunners bool, ctx context.Context) (string, string, error) {
	var err error
	if installationID == 0 {
		if c.cfg.UseOrgRunners {
			installationID, err = c.getInstallationID("org", owner, repository, ctx)
		} else {
			installationID, err = c.getInstallationID("repo", owner, repository, ctx)
		}
		if err != nil {
			slog.Error("Failed to get installation ID", slog.String("error", err.Error()))
			return "", "", err
		}
	}

	client, err := c.getClient(installationID)
	if err != nil {
		return "", "", err
	}

	var githubRunnerConfig string
	var useJitConfigStr string
	if c.cfg.Ephemeral && c.cfg.UseJitConfig {
		useJitConfigStr = "true"
		githubRunnerConfig, err = c.generateJitConfig(owner, repository, instanceName, useOrgRunners, ctx, client)
		if err != nil {
			slog.Error("Failed to generate JIT config", slog.String("error", err.Error()))
			return "", "", err
		}
	} else {
		useJitConfigStr = "false"
		githubRunnerConfig, err = c.generateStandardConfig(owner, repository, useOrgRunners, ctx, client)
		if err != nil {
			slog.Error("Failed to generate standard config", slog.String("error", err.Error()))
			return "", "", err
		}
	}
	return githubRunnerConfig, useJitConfigStr, nil
}

func (c *Client) RemoveRunnerForInstance(instance *compute.Instance, ctx context.Context) (bool, error) {
	var installationId int64

	repo := instance.Labels["ghr-repo"]
	owner := instance.Labels["ghr-owner"]
	runnerType := instance.Labels["ghr-type"]

	installationId, err := c.getInstallationID(runnerType, owner, repo, ctx)
	if err != nil {
		err = fmt.Errorf("failed to get installation id: %w", err)
		return false, err
	}

	client, err := c.getClient(installationId)
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

func (c *Client) getInstallationID(runnerType, owner, repo string, ctx context.Context) (int64, error) {
	switch runnerType {
	case "repo":
		installation, _, err := c.appsClient.Apps.FindRepositoryInstallation(ctx, owner, repo)
		if err != nil {
			err = fmt.Errorf("failed to find installation: %w", err)
			return 0, err
		}
		return installation.GetID(), nil
	case "org":
		installation, _, err := c.appsClient.Apps.FindOrganizationInstallation(ctx, owner)
		if err != nil {
			err = fmt.Errorf("failed to find installation: %w", err)
			return 0, err
		}
		return installation.GetID(), nil
	}
	return 0, fmt.Errorf("invalid runner type: %s", runnerType)
}

func (c *Client) GetJobStatus(jobID, installationID int64, owner, repo string, ctx context.Context) (string, error) {
	var err error
	if installationID == 0 {
		if c.cfg.UseOrgRunners {
			installationID, err = c.getInstallationID("org", owner, repo, ctx)
		} else {
			installationID, err = c.getInstallationID("repo", owner, repo, ctx)
		}
		if err != nil {
			slog.Error("Failed to get installation ID", slog.String("error", err.Error()))
			return "", err
		}
	}

	client, err := c.getClient(installationID)
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
