package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v60/github"
	"google.golang.org/api/compute/v1"
)

func (h *ControlPlaneHandler) StopRunner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Error(fmt.Sprintf("Invalid request method: %s", r.Method))
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()

	slog.Debug(fmt.Sprintf("Creating apps client for app %d", h.cfg.AppID))

	slog.Debug("Listing managed instances")

	instanceList, err := h.computeService.Instances.AggregatedList(h.cfg.ProjectID).Filter("labels.ghr-managed=true").Do()
	if err != nil {
		slog.Error(fmt.Sprintf("Error listing instances: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Debug(fmt.Sprintf("Found %d instances", len(instanceList.Items)))
	var wg sync.WaitGroup
	for _, zone := range instanceList.Items {
		for _, instance := range zone.Instances {
			wg.Add(1)
			go h.processInstance(instance, &wg, ctx)
		}
	}
	slog.Info("Waiting for all instances to be processed")
	wg.Wait()
	slog.Info("All instances processed. Exiting.")

	fmt.Fprint(w, "Success!")
}

func (h *ControlPlaneHandler) processInstance(instance *compute.Instance, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	creationTimestamp, err := time.Parse(time.RFC3339, instance.CreationTimestamp)
	if err != nil {
		slog.Error(fmt.Sprintf("Error parsing creation timestamp: %s", err))
		return
	}
	if creationTimestamp.Add(5 * time.Minute).After(time.Now()) {
		slog.Info(fmt.Sprintf("Instance %s has not been running for 5 minutes, skipping", instance.Name))
		return
	}
	repo := instance.Labels["ghr-repo"]
	owner := instance.Labels["ghr-owner"]
	runnerType := instance.Labels["ghr-type"]
	var installationId int64
	switch runnerType {
	case "repo":
		installation, _, err := h.appsClient.Apps.FindRepositoryInstallation(ctx, owner, repo)
		if err != nil {
			slog.Error(fmt.Sprintf("Error finding installation: %s", err))
			return
		}
		installationId = installation.GetID()
	case "org":
		installation, _, err := h.appsClient.Apps.FindOrganizationInstallation(ctx, owner)
		if err != nil {
			slog.Error(fmt.Sprintf("Error finding installation: %s", err))
			return
		}
		installationId = installation.GetID()
	}
	itr, err := ghinstallation.New(http.DefaultTransport, h.cfg.AppID, installationId, h.privateKeyBytes)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating installation client: %s", err))
		return
	}
	ghClient := github.NewClient(&http.Client{Transport: itr})

	var runners *github.Runners
	switch runnerType {
	case "repo":
		runners, _, err = ghClient.Actions.ListRunners(ctx, owner, repo, nil)
		if err != nil {
			slog.Error(fmt.Sprintf("Error listing runners: %s", err))
			return
		}
	case "org":
		runners, _, err = ghClient.Actions.ListOrganizationRunners(ctx, owner, nil)
		if err != nil {
			slog.Error(fmt.Sprintf("Error listing runners: %s", err))
			return
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
		return
	}
	if runner.GetBusy() {
		slog.Info(fmt.Sprintf("Runner %s is busy, skipping", instance.Name))
		return
	}
	switch runnerType {
	case "repo":
		_, err = ghClient.Actions.RemoveRunner(ctx, owner, repo, runner.GetID())
		if err != nil {
			slog.Error(fmt.Sprintf("Error removing runner: %s", err))
			return
		}
	case "org":
		_, err = ghClient.Actions.RemoveOrganizationRunner(ctx, owner, runner.GetID())
		if err != nil {
			slog.Error(fmt.Sprintf("Error removing runner: %s", err))
			return
		}
	}
	slog.Info(fmt.Sprintf("Runner %s removed", instance.Name))
	slog.Info(fmt.Sprintf("Deleting instance %s", instance.Name))
	zoneSplit := strings.Split(instance.Zone, "/")
	zone := zoneSplit[len(zoneSplit)-1]
	op, err := h.computeService.Instances.Delete(h.cfg.ProjectID, zone, instance.Name).Do()
	if err != nil {
		slog.Error(fmt.Sprintf("Error deleting instance: %s", err))
		return
	}
	for op.Status != "DONE" {
		time.Sleep(1 * time.Second)
		op, err = h.computeService.ZoneOperations.Get(h.cfg.ProjectID, zone, op.Name).Do()
		if err != nil {
			slog.Error(fmt.Sprintf("Error getting operation status: %s", err))
			return
		}
	}
	if op.Error != nil {
		errorMessages := make([]string, 0, len(op.Error.Errors))
		for _, e := range op.Error.Errors {
			errorMessages = append(errorMessages, e.Message)
		}
		slog.Error(fmt.Sprintf("Error deleting instance: %s", strings.Join(errorMessages, ";")))
		return
	}
	slog.Info(fmt.Sprintf("Instance %s deleted", instance.Name))
}
