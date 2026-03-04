package auth

import (
	"errors"
	"net/mail"
	"strings"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *AuthRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)

	if r.Email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email format")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
