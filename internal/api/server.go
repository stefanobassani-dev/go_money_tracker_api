package api

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stefanobassani-dev/money-tracker/internal/auth"
	"github.com/stefanobassani-dev/money-tracker/internal/config"
	"github.com/stefanobassani-dev/money-tracker/internal/handler"
	"github.com/stefanobassani-dev/money-tracker/internal/service"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	authHandler := setupAuth()
	r.Mount("/auth", authHandler.Routes())

	return r
}

func (s *Server) Run() {
	port := ":" + s.cfg.Server.Port

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(port, s.mount()); err != nil {
		log.Fatal(err)
	}
}

func setupAuth() *handler.AuthHandler {
	provider := auth.NewEmailPasswordAuth()
	authService := service.NewService(provider)
	authHandler := handler.NewAuthHandler(authService)

	return authHandler
}
