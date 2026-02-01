package account

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
	accountService := setup(ctx)

	credentialID := "6b48fb496b4f41a1bc31f1aa5264e78d"
	userID := "ccf01d88-ece9-4d47-a0b6-ab3973f7eb66"

	accountService.SaveAccountsByCredentialID(ctx, credentialID, userID)

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
	accountRepo := postgres.NewAccountRepository(pool)
	return NewService(accountRepo, tinkClient)
}
