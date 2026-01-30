package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
)

func NewClient(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}

var (
	CredentialsJobsQueue = "credential_ingestion"
	SyncJobsQueue        = "data_sync"
)
