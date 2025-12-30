package tink

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	apiJson "github.com/stefanobassani-dev/money-tracker/internal/api/json"
	tinkapi "github.com/stefanobassani-dev/money-tracker/internal/integration/tinkapi"
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

	r.Post("/webhooks", h.CreateWebhook)

	return r
}

func (h *Handler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	var request tinkapi.WebhookEndpoint
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apiJson.Error(w, http.StatusBadRequest, err.Error())
	}

	ctx := r.Context()
	webhook, err := h.service.CreateWebhook(ctx, request)
	if err != nil {
		apiJson.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	apiJson.Success(w, http.StatusCreated, webhook)
}
