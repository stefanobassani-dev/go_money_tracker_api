package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) LinkTinkUser(ctx context.Context, userID string, tinkUserID string) error {
	query := `
        UPDATE users 
        SET tink_user_id = $1, updated_at = NOW() 
        WHERE user_id = $2
    `

	result, err := r.db.Exec(ctx, query, tinkUserID, userID)
	if err != nil {
		return fmt.Errorf("errore durante l'update del tink_user_id: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("nessun utente trovato con ID %s", userID)
	}

	return nil
}

func (r *Repository) SaveTinkCredential(ctx context.Context, userID string, credID string) error {
	query := `
        INSERT INTO tink_credentials (user_id, credentials_id, updated_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (credentials_id) 
        DO UPDATE SET 
            updated_at = NOW(),
            status = 'UPDATED';
    `

	_, err := r.db.Exec(ctx, query, userID, credID)
	if err != nil {
		return fmt.Errorf("database error saving credential: %w", err)
	}
	return nil
}

func (r *Repository) GetTinkUserID(ctx context.Context, externalUserID string) (string, error) {
	query := `
		SELECT user_id 
		FROM users 
		WHERE external_user_id = $1;
	`

	var tinkID string
	err := r.db.QueryRow(ctx, query, externalUserID).Scan(&tinkID)
	if err != nil {
		return "", fmt.Errorf("database error retrieving user: %w", err)
	}

	return tinkID, nil
}
