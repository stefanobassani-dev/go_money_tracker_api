package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	customMiddleware "github.com/stefanobassani-dev/money-tracker/internal/api/middleware"
)

type Server struct {
	app *App
}

func NewServer(app *App) *Server {
	return &Server{
		app: app,
	}
}

func (s *Server) Mount() http.Handler {
	r := chi.NewRouter()
	setupMiddleware(r)

	r.Mount("/auth", s.app.authHandler.Routes())
	r.Mount("/tink", s.app.tinkHandler.Routes())

	return r
}

func (s *Server) Run() {
	port := ":" + s.app.cfg.Server.Port

	slog.Info("starting application", "port", port, "env", "dev")
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
