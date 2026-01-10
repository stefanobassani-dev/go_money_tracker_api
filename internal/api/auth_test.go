package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stefanobassani-dev/money-tracker/internal/auth/jwt"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
	"golang.org/x/crypto/bcrypt"
)

var (
	testApp http.Handler
	db      *pgxpool.Pool
)

func TestMain(m *testing.M) {
	cfg := config.Load()
	ctx := context.Background()

	var err error
	db, err = pgxpool.New(ctx, cfg.DB.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.Expiry)
	s := NewServer(cfg, db, jwtManager)

	testApp = s.Mount()
	code := m.Run()

	db.Close()
	os.Exit(code)
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	db.Exec(ctx, "DELETE from users")
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), 12)
	if err != nil {
		log.Fatal(err)
	}
	db.Exec(ctx, "INSERT INTO USERS (emial, password) VALUES ('stefano@gmail.com', $1)", hash)

	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
	}{
		{
			name: "Successo con credenziali corrette",
			body: map[string]string{
				"email":    "stefano@gmail.com",
				"password": "password",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Errore: Password sbagliata",
			body: map[string]string{
				"email":    "stefano@gmail.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Errore: Utente non esistente",
			body: map[string]string{
				"email":    "anonimo@gmail.com",
				"password": "password123",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Errore: Body malformato (manca password)",
			body: map[string]string{
				"email": "stefano@gmail.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Errore: Body malformato (manca email)",
			body: map[string]string{
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Errore: Body malformato (email invalida)",
			body: map[string]string{
				"email":    "stefano_gmail.com",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			testApp.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Caso '%s': atteso %d, ottenuto %d. Body: %s",
					tt.name, tt.expectedStatus, w.Code, w.Body.String())
			}

		})
	}

}
