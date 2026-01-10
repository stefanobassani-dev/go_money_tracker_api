package itest

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	testApp http.Handler
	db      *pgxpool.Pool
)

type Test struct {
	name           string
	body           map[string]string
	expectedStatus int
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	flushDBAndCreateUser(ctx)

	tests := []Test{
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

	runTests(t, "/auth/login", tests)
}

func TestRegister(t *testing.T) {
	ctx := context.Background()
	flushDBAndCreateUser(ctx)

	tests := []Test{
		{
			name: "Successo con utente non esistente",
			body: map[string]string{
				"email":    "johndoe@gmail.com",
				"password": "password",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Errore: Utente già esistente",
			body: map[string]string{
				"email":    "stefano@gmail.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusConflict,
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

	runTests(t, "/auth/register", tests)
}

func flushDBAndCreateUser(ctx context.Context) {
	db.Exec(ctx, "DELETE from users")
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), 12)
	if err != nil {
		log.Fatal(err)
	}
	db.Exec(ctx, "INSERT INTO USERS (email, password) VALUES ('stefano@gmail.com', $1)", hash)
}

func runTests(t *testing.T, path string, tests []Test) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)

			req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(jsonBody))
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
