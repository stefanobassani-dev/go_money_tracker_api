package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type Queue struct {
	rdb *redis.Client
}

func NewQueue(rdb *redis.Client) *Queue {
	return &Queue{rdb: rdb}
}

func (q *Queue) EnqueueCredential(ctx context.Context, credID, userID string) error {
	payload := domain.CredentialJobPayload{
		CredentialID: credID,
		UserID:       userID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	err = q.rdb.LPush(ctx, CredentialsJobsQueue, data).Err()
	if err != nil {
		slog.Error("failed to push to redis", "credential_id", credID, "user_id", userID)
		return fmt.Errorf("failed to push to redis: %w", err)
	}

	return nil
}

func (q *Queue) DequeueCredential(ctx context.Context, timeout time.Duration) (domain.CredentialJobPayload, error) {
	result, err := q.rdb.BRPop(ctx, timeout, CredentialsJobsQueue).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.CredentialJobPayload{}, domain.ErrNoJob
		}
		return domain.CredentialJobPayload{}, fmt.Errorf("failed to dequeue: %w", err)
	}

	if len(result) != 2 {
		return domain.CredentialJobPayload{}, fmt.Errorf("invalid result from redis: %v", result)
	}

	var payload domain.CredentialJobPayload
	if err := json.Unmarshal([]byte(result[1]), &payload); err != nil {
		return domain.CredentialJobPayload{}, fmt.Errorf("failed to unmarshal job payload: %w", err)
	}

	return payload, nil
}
