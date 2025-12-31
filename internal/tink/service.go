package tink

import (
	"context"

	"github.com/stefanobassani-dev/money-tracker/internal/integration/tinkapi"
)

type Service struct {
	tinkClient   *tinkapi.Client
	tokenManager *tinkapi.TokenManager
	repo         *Repository
}

func NewService(tinkClient *tinkapi.Client, tokenManager *tinkapi.TokenManager, repo *Repository) *Service {
	return &Service{
		tinkClient:   tinkClient,
		tokenManager: tokenManager,
		repo:         repo,
	}
}

func (s *Service) CreateWebhook(ctx context.Context, request tinkapi.WebhookEndpoint) (tinkapi.WebhookEndpoint, error) {
	clientToken, err := s.tokenManager.GetToken(ctx)
	if err != nil {
		return tinkapi.WebhookEndpoint{}, err
	}
	webhook, err := s.tinkClient.CreateWebhook(ctx, request.URL, request.Description, request.EnabledEvents, clientToken)
	if err != nil {
		return tinkapi.WebhookEndpoint{}, err
	}

	err = s.repo.SaveWebhook(ctx, webhook)
	if err != nil {
		return tinkapi.WebhookEndpoint{}, err
	}
	return webhook, nil
}

func (s *Service) DeleteWebhook(ctx context.Context, webhookID string) error {
	clientToken, err := s.tokenManager.GetToken(ctx)
	if err != nil {
		return err
	}

	err = s.tinkClient.DeleteWebhook(ctx, webhookID, clientToken)
	if err != nil {
		return err
	}

	return s.repo.DeleteWebhook(ctx, webhookID)
}

func (s *Service) ListCredentials(ctx context.Context, externalUserID string) (tinkapi.CredentialResponse, error) {
	clientToken, err := s.tokenManager.GetToken(ctx)
	if err != nil {
		return tinkapi.CredentialResponse{}, err
	}

	code, err := s.tinkClient.AuthorizationGrant(ctx, externalUserID, clientToken)
	if err != nil {
		return tinkapi.CredentialResponse{}, err
	}

	tokenRes, err := s.tinkClient.GetUserAccessToken(ctx, code)
	if err != nil {
		return tinkapi.CredentialResponse{}, err
	}

	return s.tinkClient.ListCredentials(ctx, tokenRes.AccessToken)
}
