package domain

import "context"

type AuthProvider interface {
	Authenticate(ctx context.Context, email string, password string) error
}
