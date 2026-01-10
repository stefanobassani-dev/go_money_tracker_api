package domain

import "context"

type AuthProvider interface {
	VerifyOrCreate(ctx context.Context, email string, password string) (*User, error)
}
