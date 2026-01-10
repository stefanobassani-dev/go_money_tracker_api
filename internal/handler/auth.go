package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/api/json"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
	"github.com/stefanobassani-dev/money-tracker/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	var req LoginRequest
	err := json.Decode(r, &req)
	if err != nil {
		json.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	user, token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			json.Error(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		json.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Success(w, http.StatusOK, map[string]any{"user": user, "token": token})
}
