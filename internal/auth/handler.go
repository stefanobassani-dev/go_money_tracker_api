package auth

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/api/json"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type Handler struct {
	service domain.AuthService
}

func NewAuthHandler(service domain.AuthService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", h.Login)
	r.Post("/register", h.Register)

	return r
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.Decode(r, &req)
	if err != nil {
		json.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if err = req.Validate(); err != nil {
		json.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			json.Error(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}
		if errors.Is(err, domain.ErrUserNotFound) {
			json.Error(w, http.StatusUnauthorized, "User not found")
			return
		}

		json.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Success(w, http.StatusOK, map[string]any{"user": user, "token": token})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.Decode(r, &req)
	if err != nil {
		json.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if err = req.Validate(); err != nil {
		json.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			json.Error(w, http.StatusConflict, "User already exists")
			return
		}
		json.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Success(w, http.StatusCreated, nil)
}
