package domain

import (
	"context"
)

type TinkClient interface {
	GetUserToken(ctx context.Context, code string) (string, error)
	CreateUser(ctx context.Context, externalID string) (string, error)
	GetUserByExternalID(ctx context.Context, externalID string) (string, error)
	GetAuthorizationGrant(ctx context.Context, externalID string, scopes []string) (string, error)
	ExchangeUserToken(ctx context.Context, externalID string, scopes []string) (string, error)
	GetAuthorizationGrantDelegate(ctx context.Context, externalID string, scopes []string) (string, error)
	BuildAuthURL(code string, state string) string
	GetUserCredential(ctx context.Context, credentialID string, externalID string) (Credential, error)
	ProviderConsent(ctx context.Context, externalUserID string) ([]ProviderConsent, error)
	ListAccounts(ctx context.Context, externalUserID string) ([]Account, error)
	FetchTransactions(ctx context.Context, externalUserID string, saveTransactions func(context.Context, []Transaction) error) error
}
