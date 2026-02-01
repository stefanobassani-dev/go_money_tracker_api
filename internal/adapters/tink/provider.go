package tink

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

func (c *Client) ProviderConsent(ctx context.Context, externalUserID string) ([]domain.ProviderConsent, error) {
	scopes := []string{"credentials:refresh", "provider-consents:read"}
	accessToken, err := c.ExchangeUserToken(ctx, externalUserID, scopes)
	if err != nil {
		slog.Error("failed to exchange tink user access token")
		return []domain.ProviderConsent{}, err
	}

	var res = struct {
		ProviderConsents []ProviderConsent
	}{}
	props := Props{
		ctx:         ctx,
		httpClient:  c.httpClient,
		baseUrl:     c.Cfg.BaseUrl,
		path:        string(PathProviderConsent),
		method:      http.MethodGet,
		contentType: ContentTypeJSON,
		body:        nil,
		result:      &res,
		token:       accessToken,
	}

	err = call(props)
	if err != nil {
		return []domain.ProviderConsent{}, err
	}

	return ToDomainProviderConsentList(res.ProviderConsents), nil
}
