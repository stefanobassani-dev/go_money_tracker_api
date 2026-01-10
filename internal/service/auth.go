package service

import (
	"context"

	"github.com/stefanobassani-dev/money-tracker/internal/auth/jwt"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type AuthService struct {
	provider domain.AuthProvider
	jwt      *jwt.Manager
}

func NewService(provider domain.AuthProvider, jwt *jwt.Manager) *AuthService {
	return &AuthService{
		provider: provider,
		jwt:      jwt,
	}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*domain.User, *string, error) {
	user, err := s.provider.VerifyOrCreate(ctx, email, password)
	if err != nil {
		return nil, nil, err
	}

	token, err := s.jwt.Generate(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, &token, nil
}
