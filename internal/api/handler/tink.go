package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/stefanobassani-dev/money-tracker/internal/service"
)

type TinkHandler struct {
	service *service.TinkService
}

func NewTinkHandler(service *service.TinkService) *TinkHandler {
	return &TinkHandler{
		service: service,
	}
}

func (h *TinkHandler) Routes() chi.Router {
	r := chi.NewRouter()

	return r
}
