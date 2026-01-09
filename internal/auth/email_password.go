package auth

import (
	"context"
	"log"
)

type EmailPasswordAuth struct {
}

func NewEmailPasswordAuth() *EmailPasswordAuth {
	return &EmailPasswordAuth{}
}

func (auth *EmailPasswordAuth) Authenticate(context context.Context, username string, password string) error {
	log.Printf("Authenticate: username=%s, password=%s", username, password)
	return nil
}
