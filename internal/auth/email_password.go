package auth

import "context"

type EmailPasswordAuth struct {
}

func (auth *EmailPasswordAuth) Authenticate(context context.Context, username string, password string) error {
	return nil
}
