package auth

import (
	"log"
	"net/http"

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

	r.Post("/register", h.register)
	r.Post("/delete", h.deleteUser)
	r.Post("/callback", h.callback)

	return r
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	h.service.register()
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	h.service.deleteUser()
}

func (h *Handler) callback(w http.ResponseWriter, r *http.Request) {
	log.Println(w)
}
