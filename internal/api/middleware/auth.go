package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/stefanobassani-dev/money-tracker/internal/api/json"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(tokenService domain.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			userID, err := tokenService.Validate(tokenString)
			if err != nil {
				json.Error(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
