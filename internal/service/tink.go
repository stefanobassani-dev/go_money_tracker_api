package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
	"github.com/stefanobassani-dev/money-tracker/internal/tink"
)

type TinkService struct {
	tink           domain.TinkClient
	userRepo       domain.UserRepository
	credentialRepo domain.CredentialRepository
}

func NewTinkService(t domain.TinkClient, userRepo domain.UserRepository,
	credentialRepo domain.CredentialRepository) *TinkService {
	return &TinkService{
		tink:           t,
		userRepo:       userRepo,
		credentialRepo: credentialRepo,
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
		var tErr *tink.TinkError
		if errors.As(err, &tErr) && tErr.StatusCode == http.StatusConflict {
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
		return err
	}

	return s.credentialRepo.CreateCredential(ctx, credential)
}
