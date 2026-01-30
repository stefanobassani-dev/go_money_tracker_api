package domain

import (
	"context"
	"time"
)

type Credential struct {
	CredentialID string
	TinkUserID   string
	UserID       string
	ProviderName string
	Status       string
	ExpiresAt    time.Time
	LastUpdated  time.Time
	Type         string
}

type CredentialRepository interface {
	CreateCredential(ctx context.Context, credential Credential) error
	CreatePendingCredential(ctx context.Context, credentialID string, userID string) error
}
