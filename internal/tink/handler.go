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
	r.Delete("/webhooks/{id}", h.DeleteWebhook)
	r.Get("/credentials", h.ListCredentials)

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

func (h *Handler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	webhookID := chi.URLParam(r, "id")
	if webhookID == "" {
		apiJson.Error(w, http.StatusBadRequest, "webhook id required")
	}

	ctx := r.Context()
	err := h.service.DeleteWebhook(ctx, webhookID)
	if err != nil {
		apiJson.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	apiJson.Success(w, http.StatusOK, nil)
}

func (h *Handler) ListCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	credentialResponse, err := h.service.ListCredentials(ctx, "2d7b9b46-94fe-435e-aa7d-95ac55fc188d")
	if err != nil {
		apiJson.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	apiJson.Success(w, http.StatusOK, credentialResponse)
}
