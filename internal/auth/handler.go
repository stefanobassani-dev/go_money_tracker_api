package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	apiJson "github.com/stefanobassani-dev/money-tracker/internal/api/json"
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
	r.Get("/callback", h.callback)

	return r
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	url, err := h.service.register(ctx)
	if err != nil {
		apiJson.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	apiJson.Success(w, http.StatusOK, url)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.service.deleteUser(ctx)
}

func (h *Handler) callback(w http.ResponseWriter, r *http.Request) {
	err := h.service.ProcessCallback(r)
	if err != nil {
		apiJson.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	apiJson.Success(w, http.StatusOK, nil)
}
