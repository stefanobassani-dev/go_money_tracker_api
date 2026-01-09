package service

import (
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type AuthService struct {
	provider domain.AuthProvider
}

func NewService(provider domain.AuthProvider) *AuthService {
	return &AuthService{
		provider: provider,
	}
}

func (s *AuthService) Login() error {
	s.provider.Authenticate(nil, "", "")
	return nil
}
