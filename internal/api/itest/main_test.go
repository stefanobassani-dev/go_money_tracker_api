package itest

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
	"github.com/stefanobassani-dev/money-tracker/internal/api"
	"github.com/stefanobassani-dev/money-tracker/internal/auth"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
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

func TestMain(m *testing.M) {
	cfg := config.Load()
	ctx := context.Background()

	var err error
	db, err = pgxpool.New(ctx, cfg.DB.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	jwtManager := auth.NewManager(cfg.JWT.Secret, cfg.JWT.Expiry)
	s := api.NewServer(cfg, db, jwtManager)

	testApp = s.Mount()
	code := m.Run()

	db.Close()
	os.Exit(code)
}

func runTests(t *testing.T, path string, method string, tests []Test) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)

			req := httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
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
