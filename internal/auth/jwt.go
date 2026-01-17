package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type Manager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewManager(secretKey string, duration time.Duration) *Manager {
	return &Manager{secretKey, duration}
}

func (m *Manager) Generate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(m.tokenDuration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.secretKey))
}

func (m *Manager) Validate(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.secretKey), nil
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
