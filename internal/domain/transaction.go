package domain

import (
	"context"
	"time"
)

type Transaction struct {
	ID                    string `json:"id"`
	AccountID             string `json:"account_id"`
	ProviderTransactionID string `json:"provider_transaction_id"`

	// Amount
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`

	// Descriptions
	Description    string `json:"description"`     // Display description
	RawDescription string `json:"raw_description"` // Original description

	// Dates
	Date       time.Time `json:"date"`        // Transaction date
	BookedDate time.Time `json:"booked_date"` // Booked date
	ValueDate  time.Time `json:"value_date"`  // Value date

	// Status
	Status string `json:"status"` // PENDING, BOOKED

	// Categorization
	CategoryID string `json:"category_id"`

	// Merchant info (opzionale)
	MerchantName         string `json:"merchant_name,omitempty"`
	MerchantCategoryCode string `json:"merchant_category_code,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TransactionRepository interface {
	UpsertBatch(ctx context.Context, txs []Transaction) error
}

type TransactionService interface {
	SaveTransactions(ctx context.Context, userID string) error
}
