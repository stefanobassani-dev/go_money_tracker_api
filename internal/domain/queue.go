package domain

import (
	"context"
	"time"
)

type CredentialJobPayload struct {
	CredentialID string
	UserID       string
}

type SyncJobPayload struct {
	UserID       string
	CredentialID string
}

type TransactionJobPayload struct {
	UserID    string
	AccountID string
}

type Queue interface {
	EnqueueCredential(ctx context.Context, credID, userID string) error
	DequeueCredential(ctx context.Context, timeout time.Duration) (CredentialJobPayload, error)
	EnqueueSync(ctx context.Context, credID, userID string) error
	DequeueSync(ctx context.Context, timeout time.Duration) (SyncJobPayload, error)
}
