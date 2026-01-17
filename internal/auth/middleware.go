package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/stefanobassani-dev/money-tracker/internal/api/json"
)

type contextKey string

const UserIDKey contextKey = "userID"

func (m *Manager) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			json.Error(w, http.StatusUnauthorized, "Authorization header missing")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			json.Error(w, http.StatusUnauthorized, "Invalid auth format")
			return
		}

		tokenString := parts[1]

		userID, err := m.Validate(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
