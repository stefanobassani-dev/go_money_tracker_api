package auth

import (
	"context"
	"errors"

	"github.com/stefanobassani-dev/money-tracker/internal/crypto"
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

func (p *EmailPasswordProvider) VerifyOrCreate(ctx context.Context, email string, password string) (*domain.User, error) {
	user, err := p.repo.FindByEmail(ctx, email)

	if errors.Is(err, domain.ErrUserNotFound) {
		hashedPw, err := crypto.HashPassword(password)
		if err != nil {
			return nil, err
		}

		user = &domain.User{
			Email:    email,
			Password: hashedPw,
		}

		err = p.repo.CreateUser(ctx, user)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	if !crypto.CheckPassword(password, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}
