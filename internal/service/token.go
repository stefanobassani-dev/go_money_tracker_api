package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type TokenService struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewTokenService(secretKey string, duration time.Duration) *TokenService {
	return &TokenService{secretKey, duration}
}

func (t *TokenService) Generate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(t.tokenDuration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(t.secretKey))
}

func (t *TokenService) Validate(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(j *jwt.Token) (interface{}, error) {
		if _, ok := j.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(t.secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok {
			return "", jwt.ErrTokenInvalidClaims
		}
		return userID, nil
	}

	return "", domain.ErrTokenInvalid
}
