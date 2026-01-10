package itest

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stefanobassani-dev/money-tracker/internal/api"
	"github.com/stefanobassani-dev/money-tracker/internal/auth"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
)

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
