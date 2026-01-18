package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/api/json"
	"github.com/stefanobassani-dev/money-tracker/internal/auth"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type TinkHandler struct {
	service    domain.TinkService
	jwtManager *auth.Manager
}

func NewTinkHandler(service domain.TinkService, jwtManager *auth.Manager) *TinkHandler {
	return &TinkHandler{
		service:    service,
		jwtManager: jwtManager,
	}
}

func (h *TinkHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/callback", h.HandleCallback)
	r.Group(func(r chi.Router) {
		r.Use(h.jwtManager.JWTMiddleware)

		r.Get("/link", h.GetConnectURL)
	})

	return r
}

func (h *TinkHandler) GetConnectURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok {
		json.Error(w, http.StatusUnauthorized, "unauthorized")
	}
	ctx := r.Context()

	_, err := h.service.GetOrCreateTinkUser(ctx, userID)
	if err != nil {
		json.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	link, err := h.service.GetConnectURL(ctx, userID)
	if err != nil {
		json.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := struct {
		URL string `json:"url"`
	}{
		URL: link,
	}
	json.Success(w, http.StatusOK, response)
}

func (h *TinkHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	credentialID := r.URL.Query().Get("credentialsId")
	externalUserID := r.URL.Query().Get("state")
	errorType := r.URL.Query().Get("error")
	errorDisplayMessage := r.URL.Query().Get("error_display_message")

	if errorType != "" {
		json.Error(w, http.StatusInternalServerError, errorDisplayMessage)
		return
	}

	if externalUserID == "" || credentialID == "" {
		json.Error(w, http.StatusBadRequest, "missing parameters")
		return
	}

	err := h.service.SaveCredential(r.Context(), credentialID, externalUserID)
	if err != nil {
		json.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	json.Success(w, http.StatusCreated, nil)
}
