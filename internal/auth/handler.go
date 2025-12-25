package auth

import (
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	return r
}
