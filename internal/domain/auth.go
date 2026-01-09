package domain

import "context"

type AuthService interface {
	Authenticate(ctx context.Context, email string, password string) error
}
