package auth

import (
	"fmt"
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
	url, err := h.service.register()
	if err != nil {
		apiJson.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	apiJson.Success(w, http.StatusOK, url)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	h.service.deleteUser()
}

func (h *Handler) callback(w http.ResponseWriter, r *http.Request) {
	credID := r.URL.Query().Get("credentials_id")
	if credID == "" {
		credID = r.URL.Query().Get("credentialsId")
	}

	if credID == "" {
		http.Error(w, "Credentials ID missing", http.StatusBadRequest)
		return
	}
	internalUserID := r.URL.Query().Get("state")
	fmt.Printf("Utente %s ha collegato con successo la banca. Credential ID: %s\n", internalUserID, credID)

	//TODO salvare il credentila id sulla riga dello user

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Banca collegata!</h1><p>ID Credenziale: %s</p>", credID)
}
