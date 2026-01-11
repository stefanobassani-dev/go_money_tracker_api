package service

import (
	"context"
	"errors"

	"github.com/stefanobassani-dev/money-tracker/internal/auth"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type AuthService struct {
	provider domain.AuthProvider
	jwt      *auth.Manager
	repo     domain.UserRepository
}

func NewAuthService(provider domain.AuthProvider, jwt *auth.Manager, repo domain.UserRepository) *AuthService {
	return &AuthService{
		provider: provider,
		jwt:      jwt,
		repo:     repo,
	}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*domain.User, *string, error) {
	user, err := s.provider.Authenticate(ctx, email, password)
	if err != nil {
		return nil, nil, err
	}

	token, err := s.jwt.Generate(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, &token, nil
}

func (s *AuthService) Register(ctx context.Context, email string, password string) error {
	_, err := s.repo.FindByEmail(ctx, email)

	if err == nil {
		return domain.ErrUserAlreadyExists
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    email,
		Password: hashedPassword,
	}
	return s.repo.CreateUser(ctx, user)
}
