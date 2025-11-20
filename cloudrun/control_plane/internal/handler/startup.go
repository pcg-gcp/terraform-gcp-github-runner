package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
)

func (h *ControlPlaneHandler) StartRunner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m, err := parseStartUpRequest(r)
	if err != nil {
		slog.Error(fmt.Sprintf("Error validating request: %s", err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	slog.Info(fmt.Sprintf("Proccesing event for %s/%s", m.Owner, m.Repository))

	if ok, err := h.makeStartUpDecision(m, ctx); err != nil {
		slog.Error(fmt.Sprintf("Error making startup decision: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if !ok {
		slog.Info(fmt.Sprintf("Ignoring event for %s/%s", m.Owner, m.Repository))
		_, err = fmt.Fprint(w, "Ignored")
		if err != nil {
			slog.Warn(fmt.Sprintf("Error writing response: %s", err))
		}
		return
	}

	id, err := randomHex(8)
	if err != nil {
		slog.Error(fmt.Sprintf("Error generating instance ID: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	instanceName := "ghr-" + id

	slog.Debug(fmt.Sprintf("Generating installation client for installation %d", m.InstallationID))

	if h.cfg.UseOrgRunners {
		slog.Debug(fmt.Sprintf("Creating repository registration token for %s/%s", m.Owner, m.Repository))
	} else {
		slog.Debug(fmt.Sprintf("Creating organization registration token for %s", m.Owner))
	}

	githubRunnerConfig, useJitConfigStr, err := h.githubClient.GenerateRunnerConfig(m.InstallationID, m.Owner, m.Repository, instanceName, h.cfg.UseOrgRunners, ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("Error generating Runner Config: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	runnerTypeString := "repo"
	if h.cfg.UseOrgRunners {
		runnerTypeString = "org"
	}

	err = h.gcpClient.CreateInstance(instanceName, m.Repository, m.Owner, githubRunnerConfig, useJitConfigStr, runnerTypeString, ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating instance: %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, "Success!")
	if err != nil {
		slog.Warn(fmt.Sprintf("Error writing response: %s", err))
	}
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
