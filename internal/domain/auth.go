package domain

import (
	"context"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*User, *string, error)
	Register(ctx context.Context, email, password string) error
	Authenticate(ctx context.Context, email, password string) (*User, error)
}
