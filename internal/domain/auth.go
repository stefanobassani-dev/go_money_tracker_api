package domain

import "context"

type AuthProvider interface {
	Authenticate(ctx context.Context, email string, password string) (*User, error)
}

type AuthService interface {
	Login(ctx context.Context, email string, password string) (*User, *string, error)
	Register(ctx context.Context, email string, password string) error
}
