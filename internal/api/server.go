package api

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stefanobassani-dev/money-tracker/internal/api/handler"
	"github.com/stefanobassani-dev/money-tracker/internal/auth"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
	"github.com/stefanobassani-dev/money-tracker/internal/repository/postgres"
	"github.com/stefanobassani-dev/money-tracker/internal/service"
	"github.com/stefanobassani-dev/money-tracker/internal/tink"
)

type Server struct {
	cfg        *config.Config
	db         *pgxpool.Pool
	jwtManager *auth.Manager
}

func NewServer(cfg *config.Config, db *pgxpool.Pool, jwtManager *auth.Manager) *Server {
	return &Server{
		cfg:        cfg,
		db:         db,
		jwtManager: jwtManager,
	}
}

func (s *Server) Mount() http.Handler {
	r := chi.NewRouter()
	setupMiddleware(r)

	userRepo := postgres.NewUserRepository(s.db)
	credentialRepo := postgres.NewCredentialRepository(s.db)

	authHandler := setupAuth(s, userRepo)
	r.Mount("/auth", authHandler.Routes())
	tinkHandler := setupTink(s, userRepo, credentialRepo)
	r.Mount("/tink", tinkHandler.Routes())

	return r
}

func (s *Server) Run() {
	port := ":" + s.cfg.Server.Port

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(port, s.Mount()); err != nil {
		log.Fatal(err)
	}
}

func setupMiddleware(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
}

func setupAuth(s *Server, repo *postgres.UserRepository) *handler.AuthHandler {
	provider := auth.NewEmailPasswordAuth(repo)
	authService := service.NewAuthService(provider, s.jwtManager, repo)
	authHandler := handler.NewAuthHandler(authService)

	return authHandler
}

func setupTink(s *Server, userRepo *postgres.UserRepository,
	credentialRepo *postgres.CredentialRepository) *handler.TinkHandler {
	httpClient := http.Client{
		Timeout: time.Second * 5,
	}
	tokenManager := tink.NewTokenManager(&httpClient, &s.cfg.Tink)
	tinkClient := tink.NewTinkClient(&s.cfg.Tink, tokenManager, &httpClient)
	tinkService := service.NewTinkService(tinkClient, userRepo, credentialRepo)
	tinkHandler := handler.NewTinkHandler(tinkService, s.jwtManager)
	return tinkHandler
}
