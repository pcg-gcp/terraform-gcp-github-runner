package gcp

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"google.golang.org/api/compute/v1"
)

func (c *Client) CreateInstance(instanceName, repository, owner, githubRunnerConfig, useJitConfigStr, runnerType string, ctx context.Context) error {
	slog.Debug(fmt.Sprintf("Getting instance template %s", c.cfg.InstanceTemplateName))
	template, err := c.computeService.InstanceTemplates.Get(c.cfg.ProjectID, c.cfg.InstanceTemplateName).Do()
	if err != nil {
		err = fmt.Errorf("error getting instance template: %w", err)
		slog.Error(err.Error())
		return err
	}

	slog.Debug(fmt.Sprintf("Creating config secret for instance %s", instanceName))
	secretName := fmt.Sprintf("%s-config", instanceName)
	if len(template.Properties.ServiceAccounts) == 0 {
		err = fmt.Errorf("no service account found in instance template: %s", c.cfg.InstanceTemplateName)
		slog.Error(err.Error())
		return err
	}
	serviceAccount := template.Properties.ServiceAccounts[0].Email
	err = c.addSecret(secretName, githubRunnerConfig, serviceAccount, ctx)
	if err != nil {
		err = fmt.Errorf("error creating secret: %w", err)
		slog.Error(err.Error())
		return err
	}

	zone, err := c.pickZone()
	if err != nil {
		err = fmt.Errorf("error picking random zone: %w", err)
		slog.Error(err.Error())
		return err
	}

	slog.Debug(fmt.Sprintf("Creating instance %s", instanceName))

	instance := &compute.Instance{
		Name: instanceName,
		Metadata: &compute.Metadata{
			Items: append(template.Properties.Metadata.Items, &compute.MetadataItems{
				Key:   "use_jit_config",
				Value: &useJitConfigStr,
			}, &compute.MetadataItems{
				Key:   "config_secret",
				Value: &secretName,
			}),
		},
		Labels: map[string]string{
			"ghr-managed": "true",
			"ghr-type":    runnerType,
			"ghr-repo":    strings.ToLower(repository),
			"ghr-owner":   strings.ToLower(owner),
		},
	}
	createInstanceRequest := c.computeService.Instances.Insert(c.cfg.ProjectID, zone, instance)
	createInstanceRequest = createInstanceRequest.SourceInstanceTemplate(template.SelfLink)
	op, err := createInstanceRequest.Do()
	if err != nil {
		err = fmt.Errorf("error creating instance: %w", err)
		slog.Error(err.Error())
		return err
	}
	for op.Status != "DONE" {
		time.Sleep(1 * time.Second)
		op, err = c.computeService.ZoneOperations.Get(c.cfg.ProjectID, zone, op.Name).Do()
		if err != nil {
			err = fmt.Errorf("error getting operation status: %w", err)
			slog.Error(err.Error())
			return err
		}
	}
	slog.Info(fmt.Sprintf("Instance %s created", instanceName))
	return nil
}

func (c *Client) GetInstances() (*compute.InstanceAggregatedList, error) {
	return c.computeService.Instances.AggregatedList(c.cfg.ProjectID).Filter("labels.ghr-managed=true").Do()
}

func (c *Client) DeleteInstance(instanceName, zone string, ctx context.Context) error {
	err := c.deleteSecret(fmt.Sprintf("%s-config", instanceName), ctx)
	if err != nil {
		slog.Warn(fmt.Sprintf("Error deleting secret for instance %s: %s", instanceName, err))
	}

	op, err := c.computeService.Instances.Delete(c.cfg.ProjectID, zone, instanceName).Do()
	if err != nil {
		err = fmt.Errorf("error deleting instance: %w", err)
		slog.Error(err.Error())
		return err
	}
	for op.Status != "DONE" {
		time.Sleep(1 * time.Second)
		op, err = c.computeService.ZoneOperations.Get(c.cfg.ProjectID, zone, op.Name).Do()
		if err != nil {
			err = fmt.Errorf("error getting operation status: %w", err)
			slog.Error(err.Error())
			return err
		}
	}
	if op.Error != nil {
		errorMessages := make([]string, 0, len(op.Error.Errors))
		for _, e := range op.Error.Errors {
			errorMessages = append(errorMessages, e.Message)
		}
		err = fmt.Errorf("error deleting instance: %s", strings.Join(errorMessages, ";"))
		slog.Error(err.Error())
		return err
	}
	slog.Info(fmt.Sprintf("Instance %s deleted", instanceName))
	return nil
}

func (c *Client) GetAvailableZones() ([]string, error) {
	zones, err := c.computeService.Zones.List(c.cfg.ProjectID).Filter(fmt.Sprintf("region=\"https://www.googleapis.com/compute/v1/projects/%s/regions/%s\"", c.cfg.ProjectID, c.cfg.Region)).Do()
	if err != nil {
		err = fmt.Errorf("error getting zones: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	availableZones := make([]string, 0)
	for _, zone := range zones.Items {
		if zone.Status == "UP" {
			availableZones = append(availableZones, zone.Name)
		}
	}
	return availableZones, nil
}

func (c *Client) pickZone() (string, error) {
	availableZones, err := c.GetAvailableZones()
	if err != nil {
		return "", err
	}
	if len(availableZones) == 0 {
		return "", fmt.Errorf("no available zones found")
	}

	if len(c.cfg.AllowedZones) == 0 {
		if c.cfg.UseStrictZoneOrder {
			return availableZones[0], nil
		}
		return availableZones[rand.Intn(len(availableZones))], nil
	}
	validZones := make([]string, 0, len(availableZones))
	for _, zone := range availableZones {
		for _, allowedZone := range c.cfg.AllowedZones {
			if zone == allowedZone {
				validZones = append(validZones, zone)
			}
		}
	}
	if len(validZones) == 0 {
		return "", fmt.Errorf("no valid zones found")
	}
	if c.cfg.UseStrictZoneOrder {
		return validZones[0], nil
	}
	return validZones[rand.Intn(len(validZones))], nil
}
