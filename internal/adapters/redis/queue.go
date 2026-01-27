package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type JobPayload struct {
	CredentialID string `json:"cred_id"`
	UserID       string `json:"user_id"`
}

type Queue struct {
	rdb *redis.Client
}

func NewQueue(rdb *redis.Client) *Queue {
	return &Queue{rdb: rdb}
}

func (q *Queue) EnqueueCredential(ctx context.Context, credID, userID string) error {
	payload := JobPayload{
		CredentialID: credID,
		UserID:       userID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	err = q.rdb.LPush(ctx, "credential_jobs", data).Err()
	if err != nil {
		slog.Error("failed to push to redis", "credential_id", credID, "user_id", userID)
		return fmt.Errorf("failed to push to redis: %w", err)
	}

	return nil
}
