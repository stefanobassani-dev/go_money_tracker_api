package domain

import (
	"context"
	"errors"
)

var (
	ErrInvalidAuth = errors.New("invalid auth credentials")
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*User, *string, error)
	Register(ctx context.Context, email, password string) error
	Authenticate(ctx context.Context, email, password string) (*User, error)
}
