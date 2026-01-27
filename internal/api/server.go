package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stefanobassani-dev/money-tracker/internal/api/handler"
	customMiddleware "github.com/stefanobassani-dev/money-tracker/internal/api/middleware"
	"github.com/stefanobassani-dev/money-tracker/internal/auth"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
	"github.com/stefanobassani-dev/money-tracker/internal/repository/postgres"
	"github.com/stefanobassani-dev/money-tracker/internal/service"
	"github.com/stefanobassani-dev/money-tracker/internal/tink"
)

type Server struct {
	cfg *config.Config
	db  *pgxpool.Pool
}

func NewServer(cfg *config.Config, db *pgxpool.Pool) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
	}
}

func (s *Server) Mount() http.Handler {
	r := chi.NewRouter()
	setupMiddleware(r)

	tokenService := service.NewTokenService(s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	authMw := customMiddleware.AuthMiddleware(tokenService)

	userRepo := postgres.NewUserRepository(s.db)
	credentialRepo := postgres.NewCredentialRepository(s.db)

	authHandler := setupAuth(userRepo, tokenService)
	r.Mount("/auth", authHandler.Routes())
	tinkHandler := setupTink(s, userRepo, credentialRepo, authMw)
	r.Mount("/tink", tinkHandler.Routes())

	return r
}

func (s *Server) Run() {
	port := ":" + s.cfg.Server.Port

	slog.Info("starting application", "port", s.cfg.Server.Port, "env", "dev")
	if err := http.ListenAndServe(port, s.Mount()); err != nil {
		slog.Error("server failed", "error", err)
	}
}

func setupMiddleware(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(customMiddleware.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
}

func setupAuth(repo *postgres.UserRepository,
	tokenService domain.TokenService) *handler.AuthHandler {
	provider := auth.NewEmailPasswordAuth(repo)
	authService := service.NewAuthService(provider, tokenService, repo)
	authHandler := handler.NewAuthHandler(authService)

	return authHandler
}

func setupTink(s *Server, userRepo *postgres.UserRepository,
	credentialRepo *postgres.CredentialRepository, authMw func(http.Handler) http.Handler) *handler.TinkHandler {
	httpClient := http.Client{
		Timeout: time.Second * 5,
	}
	tokenManager := tink.NewTokenManager(&httpClient, &s.cfg.Tink)
	tinkClient := tink.NewTinkClient(&s.cfg.Tink, tokenManager, &httpClient)
	tinkService := service.NewTinkService(tinkClient, userRepo, credentialRepo)
	tinkHandler := handler.NewTinkHandler(tinkService, authMw)
	return tinkHandler
}
