package domain

type TokenService interface {
	Generate(userID string) (string, error)
	Validate(tokenString string) (string, error)
}
