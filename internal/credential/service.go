package credential

import (
	"context"
	"log/slog"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type Service struct {
	credentialRepo domain.CredentialRepository
	tinkClient     domain.TinkClient
	redisQueue     domain.Queue
}

func NewService(credRepo domain.CredentialRepository, tinkClient domain.TinkClient,
	redisQueue domain.Queue) *Service {
	return &Service{
		credentialRepo: credRepo,
		tinkClient:     tinkClient,
		redisQueue:     redisQueue,
	}
}

func (s *Service) GetConnectURL(ctx context.Context, externalUserID string) (string, error) {
	scopes := []string{
		"authorization:read",
		"authorization:grant",
		"credentials:refresh",
		"credentials:read",
		"credentials:write",
		"providers:read",
		"user:read",
		"accounts:read",
		"transactions:read",
	}
	code, err := s.tinkClient.GetAuthorizationGrantDelegate(ctx, externalUserID, scopes)
	if err != nil {
		return "", err
	}

	return s.tinkClient.BuildUrl(code, externalUserID), nil
}

func (s *Service) SaveCredential(ctx context.Context, credentialID string, externalUserID string) error {
	err := s.credentialRepo.CreatePendingCredential(ctx, credentialID, externalUserID)
	if err != nil {
		slog.Error("error saving pending credential on database", "credentials_id", credentialID)
		return err
	}

	err = s.redisQueue.EnqueueCredential(ctx, credentialID, externalUserID)
	if err != nil {
		slog.Error("[CRITICAL] error enqueuing in redis", "credentials_id",
			credentialID, "user_id", externalUserID)
		// TODO Nota: Qui potremmo voler fare rollback del DB o segnare come FAILED,
		// ma per ora ritorniamo errore così il client sa che qualcosa è andato storto.
		return err
	}

	slog.Info("credential successfully enqueued", "credentials_id", credentialID,
		"user_id", externalUserID)
	return nil
}
