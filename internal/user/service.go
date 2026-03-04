package user

import (
	"context"
	"errors"
	"log/slog"

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

// GetOrCreateUser checks if a Tink user already exists for the given external ID.
// It returns the existing user ID if found; otherwise, it creates a new user
// and returns the newly generated ID.
func (s *Service) GetOrCreateTinkUser(ctx context.Context, userID string) (string, error) {
	localTinkID, err := s.userRepo.GetTinkIDByUserID(ctx, userID)
	if err == nil && localTinkID != "" {
		return localTinkID, nil
	}

	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		slog.Error("database error during user tink id retrieval", "userID", userID)
		return "", domain.ErrInternal
	}

	tinkID, err := s.tinkClient.CreateUser(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			tinkID, err = s.tinkClient.GetUserByExternalID(ctx, userID)
			if err != nil {
				if errors.Is(err, domain.ErrUserNotFound) {
					slog.Error("critical tink inconsistency", "userID", userID)
				}
				return "", domain.ErrInternal
			}
		} else {
			return "", err
		}
	}

	if err := s.userRepo.UpdateTinkID(ctx, userID, tinkID); err != nil {
		slog.Warn("asynchronous data inconsistency: tink id created but not saved locally",
			"userID", userID,
			"tinkID", tinkID,
			"error", err,
		)
	}

	return tinkID, nil
}
