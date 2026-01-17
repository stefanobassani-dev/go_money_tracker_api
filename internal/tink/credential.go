package tink

import (
	"context"
	"net/http"

	"github.com/stefanobassani-dev/money-tracker/internal/models"
)

func (c *Client) GetUserCredential(ctx context.Context, credentialID string,
	externalUserID string) (*models.Credential, error) {
	scopes := []string{"credentials:write", "credentials:read"}
	accessToken, err := c.ExchangeUserToken(ctx, externalUserID, scopes)
	if err != nil {
		return nil, err
	}

	var res models.Credential
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
		return nil, err
	}

	return &res, nil
}
