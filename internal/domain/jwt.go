package domain

import "errors"

var (
	ErrTokenInvalid          = errors.New("invalid token")
	ErrTokenInvalidSignature = errors.New("invalid token signature")
	ErrTokenExpired          = errors.New("token expired")
)

type JWTService interface {
	Generate(userID string) (string, error)
	Validate(tokenString string) (string, error)
}
