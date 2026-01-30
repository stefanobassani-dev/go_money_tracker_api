package domain

import (
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
	//Upsert(ctx context.Context, account Account) error
	//ListByUserID(ctx context.Context, userID string) ([]Account, error)
}

type AccountService interface {
}
