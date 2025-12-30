package tinkapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookEndpoint struct {
	ID            string    `json:"id,omitempty"`
	Description   string    `json:"description,omitempty"`
	URL           string    `json:"url"`
	Secret        string    `json:"secret,omitempty"`
	Disabled      bool      `json:"disabled,omitempty"`
	EnabledEvents []string  `json:"enabledEvents"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt,omitempty"`
}

func (c *Client) CreateWebhook(ctx context.Context, url string, description string,
	events []string, clientToken string) (WebhookEndpoint, error) {
	path := "/events/v2/webhook-endpoints"

	data := WebhookEndpoint{
		Description:   description,
		URL:           url,
		Disabled:      false,
		EnabledEvents: events,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return WebhookEndpoint{}, err
	}

	var res WebhookEndpoint
	err = c.call(ctx, http.MethodPost, path, "json", bytes.NewReader(body), &res, clientToken)
	if err != nil {
		return WebhookEndpoint{}, fmt.Errorf("errore creazione webhook: %w", err)
	}

	return res, nil
}

func (c *Client) DeleteWebhook(ctx context.Context, webhookID string, clientToken string) error {
	path := "/events/v2/webhook-endpoints" + "/" + webhookID

	return c.call(ctx, http.MethodDelete, path, "json", nil, nil, clientToken)
}
