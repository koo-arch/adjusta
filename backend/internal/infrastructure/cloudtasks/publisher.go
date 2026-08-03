package cloudtasks

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/oauth2/google"
)

const cloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"

type PublisherConfig struct {
	ProjectID             string
	Location              string
	Queue                 string
	HandlerURL            string
	OIDCAudience          string
	InvokerServiceAccount string
}

type Publisher struct {
	config PublisherConfig
	client *http.Client
}

func NewPublisher(ctx context.Context, cfg PublisherConfig) (*Publisher, error) {
	client, err := google.DefaultClient(ctx, cloudPlatformScope)
	if err != nil {
		return nil, err
	}
	return &Publisher{config: cfg, client: client}, nil
}

func (p *Publisher) Publish(ctx context.Context, messageID uuid.UUID) error {
	body, err := json.Marshal(map[string]string{"outbox_message_id": messageID.String()})
	if err != nil {
		return err
	}
	parent := fmt.Sprintf("projects/%s/locations/%s/queues/%s", p.config.ProjectID, p.config.Location, p.config.Queue)
	taskName := parent + "/tasks/outbox-" + messageID.String()
	requestBody := map[string]any{
		"task": map[string]any{
			"name": taskName,
			"httpRequest": map[string]any{
				"httpMethod": "POST",
				"url":        p.config.HandlerURL,
				"headers":    map[string]string{"Content-Type": "application/json"},
				"body":       base64.StdEncoding.EncodeToString(body),
				"oidcToken": map[string]string{
					"serviceAccountEmail": p.config.InvokerServiceAccount,
					"audience":            p.config.OIDCAudience,
				},
			},
		},
	}
	encoded, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	endpoint := "https://cloudtasks.googleapis.com/v2/" + url.PathEscape(parent) + "/tasks"
	endpoint = strings.ReplaceAll(endpoint, "%2F", "/")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(encoded))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusConflict {
		return nil
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return fmt.Errorf("cloud tasks create task returned %s: %s", resp.Status, strings.TrimSpace(string(responseBody)))
}
