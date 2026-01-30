package tink

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type TinkService struct {
	tink           domain.TinkClient
	userRepo       domain.UserRepository
	credentialRepo domain.CredentialRepository
	redisQueue     domain.Queue
}

func NewTinkService(t domain.TinkClient, userRepo domain.UserRepository,
	credentialRepo domain.CredentialRepository, redisQueue domain.Queue) *TinkService {
	return &TinkService{
		tink:           t,
		userRepo:       userRepo,
		credentialRepo: credentialRepo,
		redisQueue:     redisQueue,
	}
}

func (s *TinkService) GetOrCreateTinkUser(ctx context.Context, userID string) (string, error) {
	localTinkID, err := s.userRepo.GetTinkIDByUserID(ctx, userID)
	if err == nil && localTinkID != "" {
		return localTinkID, nil
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("database error during lookup: %w", err)
	}

	tinkID, err := s.tink.CreateUser(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			tinkID, err = s.tink.GetUserByExternalID(ctx, userID)
			if err != nil {
				return "", fmt.Errorf("failed to recover existing user from tink: %w", err)
			}
		} else {
			return "", fmt.Errorf("failed to create user on tink: %w", err)
		}
	}

	if err := s.userRepo.UpdateTinkID(ctx, userID, tinkID); err != nil {
		log.Printf("warning: could not save tink id %s for user %s: %v", tinkID, userID, err)
	}

	return tinkID, nil
}

func (s *TinkService) GetConnectURL(ctx context.Context, externalUserID string) (string, error) {
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
	code, err := s.tink.GetAuthorizationGrantDelegate(ctx, externalUserID, scopes)
	if err != nil {
		return "", err
	}

	return s.tink.BuildUrl(code, externalUserID), nil
}

func (s *TinkService) SaveCredential(ctx context.Context, credentialID string, externalUserID string) error {
	credential, err := s.tink.GetUserCredential(ctx, credentialID, externalUserID)
	if err != nil {
		slog.Error("error retrieving tink credential", "credentials_id", credentialID)
		return err
	}

	err = s.credentialRepo.CreateCredential(ctx, credential)
	if err == nil {
		slog.Info("credentials successfully saved on database", "credentials_id", credentialID)
	} else {
		slog.Error("error saving credentials on database", "credentials_id", credentialID)
		err := s.redisQueue.EnqueueCredential(ctx, credentialID, externalUserID)
		if err != nil {
			slog.Error("[CRITICAL] error enqueuing in redis", "credentials_id",
				credentialID, "user_id", externalUserID)
			return err
		}
	}
	return nil
}
