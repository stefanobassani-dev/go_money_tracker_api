package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
	"github.com/redis/go-redis/v9"
	"github.com/stefanobassani-dev/money-tracker/internal/account"
	"github.com/stefanobassani-dev/money-tracker/internal/adapters/postgres"
	redis2 "github.com/stefanobassani-dev/money-tracker/internal/adapters/redis"
	tink2 "github.com/stefanobassani-dev/money-tracker/internal/adapters/tink"
	customMiddleware "github.com/stefanobassani-dev/money-tracker/internal/api/middleware"
	"github.com/stefanobassani-dev/money-tracker/internal/auth"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
	"github.com/stefanobassani-dev/money-tracker/internal/credential"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
	"github.com/stefanobassani-dev/money-tracker/internal/transaction"
	"github.com/stefanobassani-dev/money-tracker/internal/user"
	"github.com/stefanobassani-dev/money-tracker/internal/worker"
)

type App struct {
	cfg *config.Config
	db  *pgxpool.Pool
	rdb *redis.Client

	//handler
	authHandler       *auth.Handler
	credentialHandler *credential.Handler

	//service
	authService        domain.AuthService
	jwtService         domain.JWTService
	credentialService  domain.CredentialService
	userService        domain.UserService
	accountService     domain.AccountService
	transactionService domain.TransactionService

	//repository
	credentialRepo    domain.CredentialRepository
	userRepo          domain.UserRepository
	accountRepository domain.AccountRepository

	redisQueue domain.Queue

	//client
	tinkClient domain.TinkClient
	httpClient *http.Client

	//worker
	CredentialWorker *worker.Credential
	SyncWorker       *worker.Sync
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
		slog.Error("failed to connect database with pgx", "error", err)
		os.Exit(1)
	}

	db, err := sqlx.Connect("pgx", cfg.DB.ConnectionString())
	if err != nil {
		slog.Error("failed to connect to database with sqlx", "error", err)
		os.Exit(1)
	}

	rdb, err := redis2.NewClient(ctx, cfg.Redis)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}

	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiry)
	redisQueue := redis2.NewQueue(rdb)
	credentialRepo := postgres.NewCredentialRepository(pool)
	userRepo := postgres.NewUserRepository(db)
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	tokenManager := tink2.NewTokenManager(httpClient, &cfg.Tink)
	tinkClient := tink2.NewTinkClient(&cfg.Tink, tokenManager, httpClient)

	userService := user.NewService(userRepo, tinkClient)

	authMiddleware := customMiddleware.AuthMiddleware(jwtService)

	credentialService := credential.NewService(credentialRepo, tinkClient, redisQueue)
	credentialHandler := credential.NewHandler(credentialService, userService, authMiddleware)

	authService := auth.NewAuthService(jwtService, userRepo)
	authHandler := auth.NewAuthHandler(authService)

	accountRepository := postgres.NewAccountRepository(pool)
	accountService := account.NewService(accountRepository, tinkClient)

	transactionRepository := postgres.NewTransactionRepository(pool)
	transactionService := transaction.NewService(transactionRepository, tinkClient)

	credentialWorker := worker.NewCredential(redisQueue, tinkClient, credentialRepo)
	syncWorker := worker.NewSync(redisQueue, accountService)

	return &App{
		cfg: cfg,
		db:  pool,
		rdb: rdb,

		authHandler:       authHandler,
		credentialHandler: credentialHandler,

		jwtService:         jwtService,
		authService:        authService,
		credentialService:  credentialService,
		userService:        userService,
		accountService:     accountService,
		transactionService: transactionService,

		userRepo:          userRepo,
		credentialRepo:    credentialRepo,
		accountRepository: accountRepository,

		tinkClient: tinkClient,
		httpClient: httpClient,

		redisQueue: redisQueue,

		CredentialWorker: credentialWorker,
		SyncWorker:       syncWorker,
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
