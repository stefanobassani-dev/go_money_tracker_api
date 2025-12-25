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
)

type Server struct {
	cfg *config.Config
	db  *pgx.Conn
}

func NewServer(cfg *config.Config, db *pgx.Conn) *Server {
	return &Server{cfg: cfg, db: db}
}

func (s *Server) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	//routes
	authHandler := auth.NewHandler(auth.NewService(s.db))
	r.Mount("/auth", authHandler.Routes())
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
