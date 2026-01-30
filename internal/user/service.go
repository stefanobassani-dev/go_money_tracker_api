package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type Service struct {
	userRepo   domain.UserRepository
	tinkClient domain.TinkClient
}

func NewService(userRepo domain.UserRepository, tinkClient domain.TinkClient) *Service {
	return &Service{
		userRepo:   userRepo,
		tinkClient: tinkClient,
	}
}

func (s *Service) GetOrCreateTinkUser(ctx context.Context, userID string) (string, error) {
	localTinkID, err := s.userRepo.GetTinkIDByUserID(ctx, userID)
	if err == nil && localTinkID != "" {
		return localTinkID, nil
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("database error during lookup: %w", err)
	}

	tinkID, err := s.tinkClient.CreateUser(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			tinkID, err = s.tinkClient.GetUserByExternalID(ctx, userID)
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
