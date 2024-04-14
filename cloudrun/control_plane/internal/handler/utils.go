package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ValidateStartUpRequest(request *http.Request) (*eventSummaryMessage, error) {
	if request.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("invalid content type %s", request.Header.Get("Content-Type"))
	}
	var message eventSummaryMessage
	if err := json.NewDecoder(request.Body).Decode(&message); err != nil {
		return nil, fmt.Errorf("error decoding request body: %s", err)
	}

	return &message, nil
}

func MakeStartUpDecision(message *eventSummaryMessage) (bool, error) {
	return true, nil
}
