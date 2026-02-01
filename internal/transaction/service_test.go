package transaction

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stefanobassani-dev/money-tracker/internal/adapters/postgres"
	tink2 "github.com/stefanobassani-dev/money-tracker/internal/adapters/tink"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
)

func TestService_SaveAccountsByCredentialID(t *testing.T) {
	ctx := context.Background()
	transactionService := setup(ctx)

	userID := "6a9f1f2e-1edc-4456-bc7b-e4900360fc40"

	transactionService.SaveTransactions(ctx, userID)

}

func setup(ctx context.Context) *Service {
	_ = godotenv.Load("../../.env")
	cfg := config.Load()
	pool, err := pgxpool.New(ctx, cfg.DB.ConnectionString())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	tokenManager := tink2.NewTokenManager(httpClient, &cfg.Tink)
	tinkClient := tink2.NewTinkClient(&cfg.Tink, tokenManager, httpClient)
	transactionRepo := postgres.NewTransactionRepository(pool)
	return NewService(transactionRepo, tinkClient)
}
