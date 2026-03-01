package auth

import (
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type JWTService struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTService(secretKey string, duration time.Duration) *JWTService {
	return &JWTService{secretKey, duration}
}

func (j *JWTService) Generate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(j.tokenDuration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		slog.Error("failed to sign JWT token", "err", err, "userID", userID)
		return "", domain.ErrInternal
	}

	return signedToken, nil
}

func (j *JWTService) Validate(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrTokenInvalidSignature
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", domain.ErrTokenExpired
		}
		return "", domain.ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok {
			return "", domain.ErrTokenInvalid
		}
		return userID, nil
	}

	return "", domain.ErrTokenInvalid
}
