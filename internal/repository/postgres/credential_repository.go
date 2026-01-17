package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type CredentialRepository struct {
	db *pgxpool.Pool
}

func NewCredentialRepository(db *pgxpool.Pool) *CredentialRepository {
	return &CredentialRepository{db: db}
}

func (r *CredentialRepository) CreateCredential(ctx context.Context, c domain.Credential) error {
	sql := `
		INSERT INTO credentials (
			tink_credential_id, 
			provider_name, 
			status, 
			updated, 
			session_expiry_date, 
			tink_user_id, 
			user_id,
		    type
		) 
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
		    $8
		)
		ON CONFLICT (tink_credential_id) 
		DO UPDATE SET 
			status = EXCLUDED.status,
			updated = EXCLUDED.updated,
			session_expiry_date = EXCLUDED.session_expiry_date,
			tink_user_id = EXCLUDED.tink_user_id;`

	_, err := r.db.Exec(ctx, sql, c.CredentialID, c.ProviderName,
		c.Status, c.LastUpdated, c.ExpiresAt, c.TinkUserID, c.UserID, c.Type)

	return err
}
