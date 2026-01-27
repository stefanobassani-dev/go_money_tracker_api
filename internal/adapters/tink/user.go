package tink

import (
	"context"
	"errors"
	"net/http"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

func (c *Client) CreateUser(ctx context.Context, externalID string) (string, error) {
	clientToken, err := c.TokenManager.GetToken(ctx)
	if err != nil {
		return "", err
	}

	req := &CreateUserRequest{
		ExternalUserID: externalID,
		Market:         c.Cfg.Market,
		Locale:         c.Cfg.Locale,
		RetentionClass: "permanent",
	}

	body, err := toJSONReader(req)
	if err != nil {
		return "", err
	}

	var res CreateUserResponse

	props := Props{
		ctx:         ctx,
		httpClient:  c.httpClient,
		baseUrl:     c.Cfg.BaseUrl,
		path:        string(PathUserCreate),
		method:      http.MethodPost,
		contentType: ContentTypeJSON,
		body:        body,
		result:      &res,
		token:       clientToken,
	}
	err = call(props)
	if err != nil {
		var tErr *TinkError
		if errors.As(err, &tErr) && tErr.StatusCode == http.StatusConflict {
			return "", domain.ErrUserAlreadyExists
		}
		return "", err
	}

	return res.UserID, nil
}

func (c *Client) GetUserByExternalID(ctx context.Context, externalID string) (string, error) {
	scopes := []string{"accounts:read"}
	accessToken, err := c.ExchangeUserToken(ctx, externalID, scopes)
	if err != nil {
		return "", err
	}

	var res User
	props := Props{
		ctx:         ctx,
		httpClient:  c.httpClient,
		baseUrl:     c.Cfg.BaseUrl,
		path:        string(PathGetUser),
		method:      http.MethodGet,
		contentType: ContentTypeJSON,
		body:        nil,
		result:      &res,
		token:       accessToken,
	}
	err = call(props)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}
