package domain

import (
	"context"

	"github.com/stefanobassani-dev/money-tracker/internal/models"
)

type TinkClient interface {
	GetUserToken(ctx context.Context, code string) (string, error)
	CreateUser(ctx context.Context, externalID string) (string, error)
	GetUserByExternalID(ctx context.Context, externalID string) (string, error)
	GetAuthorizationGrant(ctx context.Context, externalID string, scopes []string) (string, error)
	ExchangeUserToken(ctx context.Context, externalID string, scopes []string) (string, error)
	GetAuthorizationGrantDelegate(ctx context.Context, externalID string, scopes []string) (string, error)
	BuildUrl(code string, state string) string
	GetUserCredential(ctx context.Context, credentialID string, externalID string) (*models.Credential, error)
}
