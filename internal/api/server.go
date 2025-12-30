package api

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/auth"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
	tinkapi2 "github.com/stefanobassani-dev/money-tracker/internal/integration/tinkapi"
	"github.com/stefanobassani-dev/money-tracker/internal/tink"
)

type Server struct {
	cfg          *config.Config
	db           *pgx.Conn
	tinkClient   *tinkapi2.Client
	tokenManager *tinkapi2.TokenManager
}

func NewServer(cfg *config.Config, db *pgx.Conn) *Server {
	tinkClient := tinkapi2.NewTinkClient(&cfg.Tink)
	tokenManager := tinkapi2.NewTokenManager(tinkClient)
	return &Server{cfg: cfg, db: db, tinkClient: tinkClient, tokenManager: tokenManager}
}

func (s *Server) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	authHandler := setupAuth(s)
	r.Mount("/auth", authHandler.Routes())

	tinkHandler := setupTink(s)
	r.Mount("/tink", tinkHandler.Routes())

	r.Mount("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	return r
}

func (s *Server) Start() error {
	handler := s.mount()

	serverPort := s.cfg.Server.Port

	srv := &http.Server{
		Addr:         ":" + serverPort,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Println("Server is listening on port", serverPort)

	return srv.ListenAndServe()
}

func setupAuth(s *Server) *auth.Handler {
	authRepo := auth.NewRepository(s.db)
	authService := auth.NewService(authRepo, s.tinkClient, s.tokenManager)
	return auth.NewHandler(authService)
}

func setupTink(s *Server) *tink.Handler {
	tinkService := tink.NewService(s.tinkClient, s.tokenManager)
	return tink.NewHandler(tinkService)
}
