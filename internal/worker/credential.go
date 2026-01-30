package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type Credential struct {
	redisQueue     domain.Queue
	tinkClient     domain.TinkClient
	credentialRepo domain.CredentialRepository
}

func NewCredential(redisQueue domain.Queue, tinkClient domain.TinkClient,
	credentialRepo domain.CredentialRepository) *Credential {
	return &Credential{
		redisQueue:     redisQueue,
		tinkClient:     tinkClient,
		credentialRepo: credentialRepo,
	}
}

func (w *Credential) Run(ctx context.Context) {
	slog.Error("credential worker started")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			w.processNext(ctx)
		}
	}
}

func (w *Credential) processNext(ctx context.Context) {
	job, err := w.redisQueue.DequeueCredential(ctx, 5*time.Second)
	if err != nil {
		if errors.Is(err, domain.ErrNoJob) {
			return
		}
		slog.Error("redis error", "error", err)
		time.Sleep(1 * time.Second)
		return
	}

	slog.Info("processing job", "credential_id", job.CredentialID)

	credential, err := w.tinkClient.GetUserCredential(ctx, job.CredentialID, job.UserID)
	if err != nil {
		slog.Error("failed to fetch from tink", "credential_id", job.CredentialID, "error", err)
		// TODO Qui potresti implementare una logica di "Dead Letter Queue" o aggiornare lo stato a FAILED nel DB
		return
	}

	err = w.credentialRepo.CreateCredential(ctx, credential)
	if err != nil {
		slog.Error("failed to save to db", "credential_id", job.CredentialID, "error", err)
		return
	}

}
