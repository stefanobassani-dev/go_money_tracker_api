package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	GetTinkIDByUserID(ctx context.Context, userID string) (string, error)
	UpdateTinkID(ctx context.Context, userID string, newTinkID string) error
}

type UserService interface {
	GetOrCreateTinkUser(ctx context.Context, userID string) (string, error)
}

type User struct {
	ID         string
	Email      string
	Password   string
	TinkUserId *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
