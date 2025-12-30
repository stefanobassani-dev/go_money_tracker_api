package tinkapi

import (
	"bytes"
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

func (c *Client) CreateWebhook(url string, description string, events []string, clientToken string) (WebhookEndpoint, error) {
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
	err = c.call(http.MethodPost, path, "json", bytes.NewReader(body), &res, clientToken)
	if err != nil {
		return WebhookEndpoint{}, fmt.Errorf("errore creazione webhook: %w", err)
	}

	return res, nil
}
