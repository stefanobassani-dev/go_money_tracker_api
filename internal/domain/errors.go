package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var (
	ErrTokenInvalid = errors.New("invalid token")
)

var (
	ErrNoJob = errors.New("no jobs in queue")
)
