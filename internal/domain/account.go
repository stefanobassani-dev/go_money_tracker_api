package domain

import (
	"context"
	"time"
)

type AccountType string

const (
	AccountTypeChecking   AccountType = "CHECKING"
	AccountTypeSavings    AccountType = "SAVINGS"
	AccountTypeCreditCard AccountType = "CREDIT_CARD"
	AccountTypeInvestment AccountType = "INVESTMENT"
	AccountTypeLoan       AccountType = "LOAN"
	AccountTypeUnknown    AccountType = "UNKNOWN"
)

type Account struct {
	ID                     string      `json:"id"`
	UserID                 string      `json:"user_id"`
	CredentialID           *string     `json:"credential_id"`
	Name                   string      `json:"name"`
	Type                   AccountType `json:"type"`
	Balance                float64     `json:"balance"`
	Currency               string      `json:"currency"`
	FinancialInstitutionID string      `json:"financial_institution_id"`
	LastRefreshed          time.Time   `json:"last_refreshed"`
	CreatedAt              time.Time   `json:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at"`
}

type AccountRepository interface {
	CreateCredential(ctx context.Context, a Account, credentialID string) error
}

type AccountService interface {
	SaveAccountsByCredentialID(ctx context.Context, credentialID, userID string) error
}
