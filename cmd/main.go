package main

import (
	"context"
	"log"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stefanobassani-dev/money-tracker/internal/api"
	"github.com/stefanobassani-dev/money-tracker/internal/auth/jwt"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DB.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.Expiry)

	server := api.NewServer(cfg, pool, jwtManager)
	server.Run()
}
