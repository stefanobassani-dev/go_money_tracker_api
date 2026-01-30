package main

import (
	"context"

	_ "github.com/jackc/pgx/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/api"
)

func main() {
	ctx := context.Background()
	app := api.Bootstrap(ctx)
	defer app.Shutdown()

	go app.CredentialWorker.Run(ctx)
	go app.SyncWorker.Run(ctx)

	server := api.NewServer(app)
	server.Run()
}
