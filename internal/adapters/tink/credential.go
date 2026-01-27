package tink

import (
	"context"
	"net/http"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

func (c *Client) GetUserCredential(ctx context.Context, credentialID string,
	externalUserID string) (domain.Credential, error) {
	scopes := []string{"credentials:write", "credentials:read"}
	accessToken, err := c.ExchangeUserToken(ctx, externalUserID, scopes)
	if err != nil {
		return domain.Credential{}, err
	}

	var res Credential
	props := Props{
		ctx:         ctx,
		httpClient:  c.httpClient,
		baseUrl:     c.Cfg.BaseUrl,
		path:        string(PathGetCredential) + credentialID,
		method:      http.MethodGet,
		contentType: ContentTypeJSON,
		body:        nil,
		result:      &res,
		token:       accessToken,
	}

	err = call(props)
	if err != nil {
		return domain.Credential{}, err
	}

	return ToDomainCredential(&res, externalUserID), nil
}
