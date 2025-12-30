package tink

import (
	"context"

	"github.com/stefanobassani-dev/money-tracker/internal/integration/tinkapi"
)

type Service struct {
	tinkClient   *tinkapi.Client
	tokenManager *tinkapi.TokenManager
}

func NewService(tinkClient *tinkapi.Client, tokenManager *tinkapi.TokenManager) *Service {
	return &Service{
		tinkClient:   tinkClient,
		tokenManager: tokenManager,
	}
}

func (s *Service) CreateWebhook(ctx context.Context, request tinkapi.WebhookEndpoint) (tinkapi.WebhookEndpoint, error) {
	clientToken, err := s.tokenManager.GetToken()
	if err != nil {
		return tinkapi.WebhookEndpoint{}, err
	}
	return s.tinkClient.CreateWebhook(request.URL, request.Description, request.EnabledEvents, clientToken)
}
