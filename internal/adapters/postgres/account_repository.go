package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CreateCredential(ctx context.Context, a domain.Account, credentialID string) error {
	query := `
		INSERT INTO accounts (
			id, 
			user_id, 
			credential_id, 
			name, 
			type, 
			balance, 
			currency, 
			financial_institution_id, 
			last_refreshed, 
			created_at, 
			updated_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) 
		DO UPDATE SET 
			name = EXCLUDED.name,
			type = EXCLUDED.type,
			balance = EXCLUDED.balance,
			currency = EXCLUDED.currency,
			last_refreshed = EXCLUDED.last_refreshed,
			updated_at = EXCLUDED.updated_at,
			credential_id = EXCLUDED.credential_id;`

	_, err := r.db.Exec(ctx, query,
		a.ID,
		a.UserID,
		credentialID,
		a.Name,
		a.Type,
		a.Balance,
		a.Currency,
		a.FinancialInstitutionID,
		a.LastRefreshed,
		a.CreatedAt,
		a.UpdatedAt,
	)

	return err
}

func (r *AccountRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Account, error) {
	query := `
		SELECT 
			id, user_id, credential_id, name, type, balance, currency, 
			financial_institution_id, last_refreshed, created_at, updated_at
		FROM accounts
		WHERE user_id = $1`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var a domain.Account
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.CredentialID,
			&a.Name,
			&a.Type,
			&a.Balance,
			&a.Currency,
			&a.FinancialInstitutionID,
			&a.LastRefreshed,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}

	return accounts, nil
}
