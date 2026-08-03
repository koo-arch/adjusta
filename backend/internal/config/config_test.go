package config

import (
	"strings"
	"testing"
)

func TestValidateServerRequiresCompleteCloudTasksConfiguration(t *testing.T) {
	cfg := Config{
		DatabaseURL:          "postgres://example",
		SessionSecret:        "secret",
		GoogleClientID:       "client-id",
		GoogleClientSecret:   "client-secret",
		GoogleRedirectURI:    "https://example.com/callback",
		CloudTasksHandlerURL: "https://example.com/internal/tasks/google-calendar-sync",
	}

	err := cfg.validateServer()
	if err == nil {
		t.Fatal("expected incomplete Cloud Tasks configuration to fail")
	}
	if !strings.Contains(err.Error(), "CLOUD_TASKS_PROJECT_ID is not set") {
		t.Fatalf("unexpected validation error: %v", err)
	}
}
