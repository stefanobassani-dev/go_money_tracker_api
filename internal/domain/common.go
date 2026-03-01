package domain

import "errors"

var (
	ErrInternal     = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrTimeout      = errors.New("operation timed out")
)
