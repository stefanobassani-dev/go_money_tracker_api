package domain

import "errors"

var (
	ErrTokenInvalid = errors.New("invalid token")
)

type JWTService interface {
	Generate(userID string) (string, error)
	Validate(tokenString string) (string, error)
}
