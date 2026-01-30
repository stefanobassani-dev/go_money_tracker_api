package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type Sync struct {
	redisQueue domain.Queue
}

func NewSync(redisQueue domain.Queue) *Sync {
	return &Sync{
		redisQueue: redisQueue,
	}
}

func (w *Sync) Run(ctx context.Context) {
	slog.Info("sync worker started")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			w.processNext(ctx)
		}
	}
}

func (w *Sync) processNext(ctx context.Context) {
	job, err := w.redisQueue.DequeueSync(ctx, 5*time.Second)
	if err != nil {
		if errors.Is(err, domain.ErrNoJob) {
			return
		}
		slog.Error("redis error", "error", err)
		time.Sleep(1 * time.Second)
		return
	}

	slog.Info("processing sync job", "credential_id", job.CredentialID)

}
