package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"github.com/redis/go-redis/v9"
	"github.com/stefanobassani-dev/money-tracker/internal/adapters/postgres"
	redis2 "github.com/stefanobassani-dev/money-tracker/internal/adapters/redis"
	tink2 "github.com/stefanobassani-dev/money-tracker/internal/adapters/tink"
	customMiddleware "github.com/stefanobassani-dev/money-tracker/internal/api/middleware"
	"github.com/stefanobassani-dev/money-tracker/internal/auth"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
	"github.com/stefanobassani-dev/money-tracker/internal/tink"
	"github.com/stefanobassani-dev/money-tracker/internal/worker"
)

type App struct {
	cfg *config.Config
	db  *pgxpool.Pool
	rdb *redis.Client

	//handler
	authHandler *auth.AuthHandler
	tinkHandler *tink.TinkHandler

	//service
	authService  domain.AuthService
	tinkService  domain.TinkService
	tokenService domain.TokenService

	//repository
	credentialRepo domain.CredentialRepository
	userRepo       domain.UserRepository

	redisQueue domain.Queue

	//client
	tinkClient domain.TinkClient
	httpClient *http.Client

	//middleware
	authMw func(http.Handler) http.Handler

	//worker
	CredentialWorker *worker.Credential
}

func Bootstrap(ctx context.Context) *App {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.DateTime,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()

	pool, err := pgxpool.New(ctx, cfg.DB.ConnectionString())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	rdb, err := redis2.NewClient(ctx, cfg.Redis)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}

	tokenService := auth.NewTokenService(cfg.JWT.Secret, cfg.JWT.Expiry)
	authMw := customMiddleware.AuthMiddleware(tokenService)
	redisQueue := redis2.NewQueue(rdb)
	credentialRepo := postgres.NewCredentialRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	tokenManager := tink2.NewTokenManager(httpClient, &cfg.Tink)
	tinkClient := tink2.NewTinkClient(&cfg.Tink, tokenManager, httpClient)
	tinkService := tink.NewTinkService(tinkClient, userRepo, credentialRepo, redisQueue)
	tinkHandler := tink.NewTinkHandler(tinkService, authMw)

	provider := auth.NewEmailPasswordAuth(userRepo)
	authService := auth.NewAuthService(provider, tokenService, userRepo)
	authHandler := auth.NewAuthHandler(authService)

	credentialWorker := worker.NewCredential(redisQueue, tinkClient, credentialRepo)

	return &App{
		cfg: cfg,
		db:  pool,
		rdb: rdb,

		tinkHandler: tinkHandler,
		authHandler: authHandler,

		tinkService:  tinkService,
		tokenService: tokenService,
		authService:  authService,

		userRepo:       userRepo,
		credentialRepo: credentialRepo,

		tinkClient: tinkClient,
		httpClient: httpClient,

		authMw: authMw,

		redisQueue: redisQueue,

		CredentialWorker: credentialWorker,
	}
}

func (a *App) Shutdown() {
	if a.db != nil {
		a.db.Close()
		slog.Info("database connection closed")
	}
	if a.rdb != nil {
		if err := a.rdb.Close(); err != nil {
			slog.Error("failed to close redis connection", "error", err)
		} else {
			slog.Info("redis connection closed")
		}
	}
}
