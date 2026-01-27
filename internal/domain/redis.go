package domain

import "context"

type RedisQueue interface {
	EnqueueCredential(ctx context.Context, credID, userID string) error
}
