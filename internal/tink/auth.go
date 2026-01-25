package tink

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/stefanobassani-dev/money-tracker/internal/auth"
)

func (c *Client) GetAuthorizationGrant(ctx context.Context, externalID string, scopes []string) (string, error) {
	clientToken, err := c.TokenManager.GetToken(ctx)
	if err != nil {
		return "", err
	}

	req := AuthorizationRequest{
		ExternalID: externalID,
		Scope:      strings.Join(scopes, ","),
	}
	v, err := query.Values(req)
	if err != nil {
		return "", err
	}

	var res AuthorizationResponse

	props := Props{
		ctx:         ctx,
		httpClient:  c.httpClient,
		baseUrl:     c.Cfg.BaseUrl,
		path:        string(PathAuthorize),
		method:      http.MethodPost,
		contentType: ContentTypeForm,
		body:        strings.NewReader(v.Encode()),
		result:      &res,
		token:       clientToken,
	}

	err = call(props)
	if err != nil {
		return "", err
	}
	return res.Code, nil
}

func (c *Client) GetAuthorizationGrantDelegate(ctx context.Context, externalID string, scopes []string) (string, error) {
	clientToken, err := c.TokenManager.GetToken(ctx)
	if err != nil {
		return "", err
	}

	IDHint := ctx.Value(auth.UserIDKey).(string)

	req := AuthorizationRequest{
		ExternalID:    externalID,
		Scope:         strings.Join(scopes, ","),
		IDHint:        IDHint,
		ActorClientID: TinkActorClientID,
	}
	v, err := query.Values(req)
	if err != nil {
		return "", err
	}

	var res AuthorizationResponse

	props := Props{
		ctx:         ctx,
		httpClient:  c.httpClient,
		baseUrl:     c.Cfg.BaseUrl,
		path:        string(PathAuthorizeDelegate),
		method:      http.MethodPost,
		contentType: ContentTypeForm,
		body:        strings.NewReader(v.Encode()),
		result:      &res,
		token:       clientToken,
	}

	err = call(props)
	if err != nil {
		return "", err
	}
	return res.Code, nil
}

func (c *Client) GetUserToken(ctx context.Context, code string) (string, error) {
	req := TokenRequest{
		GrantType:    "authorization_code",
		ClientID:     c.Cfg.ClientId,
		ClientSecret: c.Cfg.ClientSecret,
		Code:         code,
	}
	v, err := query.Values(req)
	if err != nil {
		return "", err
	}

	var res TokenResponse

	props := Props{
		ctx:         ctx,
		httpClient:  c.httpClient,
		baseUrl:     c.Cfg.BaseUrl,
		path:        string(PathTokenExchange),
		method:      http.MethodPost,
		contentType: ContentTypeForm,
		body:        strings.NewReader(v.Encode()),
		result:      &res,
		token:       "",
	}

	err = call(props)
	if err != nil {
		return "", err
	}
	return res.AccessToken, nil
}

func (c *Client) ExchangeUserToken(ctx context.Context, externalID string, scopes []string) (string, error) {
	code, err := c.GetAuthorizationGrant(ctx, externalID, scopes)
	if err != nil {
		return "", fmt.Errorf("failed to get auth grant: %w", err)
	}

	userAccessToken, err := c.GetUserToken(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return userAccessToken, nil
}

func (c *Client) BuildUrl(code string, state string) string {
	return ApiConnectURL + "?client_id=" + c.Cfg.ClientId + "&state=" + state +
		"&redirect_uri=" + c.Cfg.RedirectUri + "&authorization_code=" + code +
		"&market=" + c.Cfg.Market + "&locale=" + c.Cfg.Locale
}
