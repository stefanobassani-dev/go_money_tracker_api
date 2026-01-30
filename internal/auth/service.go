package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type Service struct {
	provider     domain.AuthProvider
	tokenService domain.TokenService
	repo         domain.UserRepository
}

func NewAuthService(provider domain.AuthProvider, tokenService domain.TokenService, repo domain.UserRepository) *Service {
	return &Service{
		provider:     provider,
		tokenService: tokenService,
		repo:         repo,
	}
}

func (s *Service) Login(ctx context.Context, email string, password string) (*domain.User, *string, error) {
	user, err := s.provider.Authenticate(ctx, email, password)
	if err != nil {
		return nil, nil, err
	}

	token, err := s.tokenService.Generate(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, &token, nil
}

func (s *Service) Register(ctx context.Context, email string, password string) error {
	_, err := s.repo.FindByEmail(ctx, email)

	if err == nil {
		slog.Error("user already exists", "email", email)
		return domain.ErrUserAlreadyExists
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		slog.Error("database error", "err", err)
		return err
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    email,
		Password: hashedPassword,
	}
	return s.repo.CreateUser(ctx, user)
}
