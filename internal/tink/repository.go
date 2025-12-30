package tink

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/integration/tinkapi"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) SaveWebhook(ctx context.Context, webhook tinkapi.WebhookEndpoint) error {
	eventsJSON, err := json.Marshal(webhook.EnabledEvents)
	if err != nil {
		return fmt.Errorf("err marshaling events: %w", err)
	}

	query := `
        INSERT INTO tink_webhooks (webhook_id, description, url, secret, disabled, enabled_events, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (webhook_id) DO UPDATE SET
            url = EXCLUDED.url,
            secret = EXCLUDED.secret,
            enabled_events = EXCLUDED.enabled_events,
            updated_at = EXCLUDED.updated_at;
    `

	_, err = repo.db.Exec(ctx, query,
		webhook.ID,
		webhook.Description,
		webhook.URL,
		webhook.Secret,
		webhook.Disabled,
		eventsJSON,
		webhook.CreatedAt,
		webhook.UpdatedAt,
	)

	return err
}

func (repo *Repository) DeleteWebhook(ctx context.Context, id string) error {
	query := `
        DELETE FROM tink_webhooks
        WHERE webhook_id = $1;
    `

	_, err := repo.db.Exec(ctx, query, id)
	return err
}
