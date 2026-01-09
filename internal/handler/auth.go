package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", h.Login)

	return r
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.service.Login()
}
