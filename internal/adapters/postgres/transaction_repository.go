package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) UpsertBatch(ctx context.Context, txs []domain.Transaction) error {
	if len(txs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}

	query := `
		INSERT INTO transactions (
			id, 
			account_id, 
			provider_transaction_id, 
			amount, 
			currency, 
			description, 
			raw_description, 
			date, 
			booked_date, 
			value_date, 
			status, 
			category_id, 
			merchant_name, 
			merchant_category_code, 
			created_at, 
			updated_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (id) 
		DO UPDATE SET 
			amount = EXCLUDED.amount,
			currency = EXCLUDED.currency,
			description = EXCLUDED.description,
			raw_description = EXCLUDED.raw_description,
			date = EXCLUDED.date,
			booked_date = EXCLUDED.booked_date,
			value_date = EXCLUDED.value_date,
			status = EXCLUDED.status,
			category_id = EXCLUDED.category_id,
			merchant_name = EXCLUDED.merchant_name,
			merchant_category_code = EXCLUDED.merchant_category_code,
			updated_at = EXCLUDED.updated_at;`

	for _, t := range txs {
		batch.Queue(query,
			t.ID,
			t.AccountID,
			t.ProviderTransactionID,
			t.Amount,
			t.Currency,
			t.Description,
			t.RawDescription,
			t.Date,
			t.BookedDate,
			t.ValueDate,
			t.Status,
			t.CategoryID,
			t.MerchantName,
			t.MerchantCategoryCode,
			t.CreatedAt,
			t.UpdatedAt,
		)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(txs); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to upsert transaction at index %d: %w", i, err)
		}
	}

	return nil
}
