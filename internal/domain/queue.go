package domain

import (
	"context"
	"time"
)

type CredentialJobPayload struct {
	CredentialID string
	UserID       string
}

type Queue interface {
	EnqueueCredential(ctx context.Context, credID, userID string) error
	DequeueCredential(ctx context.Context, timeout time.Duration) (CredentialJobPayload, error)
	//EnqueueSyncJob()
	//DequeueSyncJob()
}
