package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	postgres2 "github.com/stefanobassani-dev/money-tracker/internal/adapters/postgres"
	tink2 "github.com/stefanobassani-dev/money-tracker/internal/adapters/tink"
	customMiddleware "github.com/stefanobassani-dev/money-tracker/internal/api/middleware"
	"github.com/stefanobassani-dev/money-tracker/internal/auth"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
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

	tokenService := auth.NewTokenService(s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	authMw := customMiddleware.AuthMiddleware(tokenService)

	userRepo := postgres2.NewUserRepository(s.db)
	credentialRepo := postgres2.NewCredentialRepository(s.db)

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

func setupAuth(repo *postgres2.UserRepository,
	tokenService domain.TokenService) *auth.AuthHandler {
	provider := auth.NewEmailPasswordAuth(repo)
	authService := auth.NewAuthService(provider, tokenService, repo)
	authHandler := auth.NewAuthHandler(authService)

	return authHandler
}

func setupTink(s *Server, userRepo *postgres2.UserRepository,
	credentialRepo *postgres2.CredentialRepository, authMw func(http.Handler) http.Handler) *tink.TinkHandler {
	httpClient := http.Client{
		Timeout: time.Second * 5,
	}
	tokenManager := tink2.NewTokenManager(&httpClient, &s.cfg.Tink)
	tinkClient := tink2.NewTinkClient(&s.cfg.Tink, tokenManager, &httpClient)
	tinkService := tink.NewTinkService(tinkClient, userRepo, credentialRepo)
	tinkHandler := tink.NewTinkHandler(tinkService, authMw)
	return tinkHandler
}
