package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	GetTinkIDByUserID(ctx context.Context, userID string) (string, error)
	UpdateTinkID(ctx context.Context, userID string, newTinkID string) error
}

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	TinkUserId *string   `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
