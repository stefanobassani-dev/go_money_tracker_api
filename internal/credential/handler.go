package credential

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/api/json"
	"github.com/stefanobassani-dev/money-tracker/internal/api/middleware"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type Handler struct {
	credentialService domain.CredentialService
	userService       domain.UserService
	authMiddleware    func(http.Handler) http.Handler
}

func NewHandler(credentialService domain.CredentialService, userService domain.UserService,
	authMiddleware func(http.Handler) http.Handler) *Handler {
	return &Handler{
		credentialService: credentialService,
		userService:       userService,
		authMiddleware:    authMiddleware,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/callback", h.HandleCallback)
	r.Group(func(r chi.Router) {
		r.Use(h.authMiddleware)

		r.Get("/link", h.GetConnectURL)
	})

	return r
}

func (h *Handler) GetConnectURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		json.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	ctx := r.Context()

	_, err := h.userService.GetOrCreateTinkUser(ctx, userID)
	if err != nil {
		json.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	link, err := h.credentialService.GetConnectURL(ctx, userID)
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

func (h *Handler) HandleCallback(w http.ResponseWriter, r *http.Request) {
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

	err := h.credentialService.SaveCredential(r.Context(), credentialID, externalUserID)
	if err != nil {
		json.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	json.Success(w, http.StatusCreated, nil)
}
