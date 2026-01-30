package auth

import (
	"context"
	"log/slog"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type EmailPasswordProvider struct {
	repo domain.UserRepository
}

func NewEmailPasswordAuth(repo domain.UserRepository) *EmailPasswordProvider {
	return &EmailPasswordProvider{
		repo: repo,
	}
}

func (p *EmailPasswordProvider) Authenticate(ctx context.Context, email string, password string) (*domain.User, error) {
	user, err := p.repo.FindByEmail(ctx, email)
	if err != nil {
		slog.Error("database error", "err", err)
		return nil, err
	}

	if !CheckPassword(password, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}
