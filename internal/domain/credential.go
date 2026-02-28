package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
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

type ProviderConsent struct {
	CredentialsID     string
	ProviderName      string
	Status            string
	SessionExpiryDate time.Time
	SessionExtendable bool
	AccountIDs        []string
	StatusUpdated     time.Time
}

type CredentialRepository interface {
	CreateCredential(ctx context.Context, credential Credential) error
	CreatePendingCredential(ctx context.Context, credentialID string, userID string) error
}

type CredentialService interface {
	GetConnectURL(ctx context.Context, externalUserID string) (string, error)
	SaveCredential(ctx context.Context, credentialID string, externalUserID string) error
}
