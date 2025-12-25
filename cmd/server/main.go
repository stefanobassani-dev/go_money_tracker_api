package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/api"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.DB.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	srv := api.NewServer(cfg, conn)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
