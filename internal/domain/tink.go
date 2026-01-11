package domain

import "context"

type TinkClient interface {
	GetClientToken(ctx context.Context, scopes []string) (string, int, error)
	//CreateUser(ctx context.Context, externalID string) (string, error)
	//GetAuthorizationGrantDelegate(ctx context.Context, externalID string, scopes []string) (string, error)
	//BuildUrl(code string) string
}
