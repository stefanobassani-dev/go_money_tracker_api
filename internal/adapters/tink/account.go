package tink

import (
	"context"
	"net/http"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

func (c *Client) ListAccounts(ctx context.Context, externalUserID string) ([]domain.Account, error) {
	scopes := []string{"accounts:read"}
	accessToken, err := c.ExchangeUserToken(ctx, externalUserID, scopes)
	if err != nil {
		return []domain.Account{}, err
	}

	var res = struct {
		accounts []Account
	}{}
	props := Props{
		ctx:         ctx,
		httpClient:  c.httpClient,
		baseUrl:     c.Cfg.BaseUrl,
		path:        string(PathListAccounts),
		method:      http.MethodGet,
		contentType: ContentTypeJSON,
		body:        nil,
		result:      &res,
		token:       accessToken,
	}

	err = call(props)
	if err != nil {
		return []domain.Account{}, err
	}

	return ToDomainAccountList(res.accounts, externalUserID), nil
}
