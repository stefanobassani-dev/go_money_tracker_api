package tink

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/stefanobassani-dev/money-tracker/internal/api/middleware"
)

type authorizationRequest struct {
	ActorClientID string `url:"actor_client_id"`
	ExternalID    string `url:"external_user_id"`
	IDHint        string `url:"id_hint,omitempty"`
	Scope         string `url:"scope"`
}

type authorizationResponse struct {
	Code string `json:"code"`
}

// GetAuthorizationGrant returns a single-use authorization code for user token generation
func (c *Client) GetAuthorizationGrant(ctx context.Context, externalID string, scopes []string) (string, error) {
	clientToken, err := c.TokenManager.GetToken(ctx)
	if err != nil {
		return "", err
	}

	req := authorizationRequest{
		ExternalID: externalID,
		Scope:      strings.Join(scopes, ","),
	}
	v, err := query.Values(req)
	if err != nil {
		return "", err
	}

	var res authorizationResponse

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

// GetAuthorizationGrantDelegate returns a a single-use code to authorizes an actor to perform operations on behalf of a specific user using a delegated grant
func (c *Client) GetAuthorizationGrantDelegate(ctx context.Context, externalID string, scopes []string) (string, error) {
	clientToken, err := c.TokenManager.GetToken(ctx)
	if err != nil {
		return "", err
	}

	IDHint := ctx.Value(middleware.UserIDKey).(string)

	req := authorizationRequest{
		ExternalID:    externalID,
		Scope:         strings.Join(scopes, ","),
		IDHint:        IDHint,
		ActorClientID: TinkActorClientID,
	}
	v, err := query.Values(req)
	if err != nil {
		return "", err
	}

	var res authorizationResponse

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

type userTokenRequest struct {
	GrantType    string `url:"grant_type"`
	ClientID     string `url:"client_id"`
	ClientSecret string `url:"client_secret"`
	Code         string `url:"code,omitempty"`
}

type userTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IdHint      string `json:"id_hint"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// GetUserToken returns an access token for user
func (c *Client) GetUserToken(ctx context.Context, code string) (string, error) {
	req := userTokenRequest{
		GrantType:    "authorization_code",
		ClientID:     c.Cfg.ClientId,
		ClientSecret: c.Cfg.ClientSecret,
		Code:         code,
	}
	v, err := query.Values(req)
	if err != nil {
		return "", err
	}

	var res userTokenResponse

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

// ExchangeUserToken calls the authorization grant method and the user token method to get a user token directly
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

// BuildAuthURL return the Tink URL to connect with a bank
func (c *Client) BuildAuthURL(code string, state string) string {
	u, _ := url.Parse(ApiConnectURL)

	q := u.Query()

	q.Set("client_id", c.Cfg.ClientId)
	q.Set("redirect_uri", c.Cfg.RedirectUri)
	q.Set("market", c.Cfg.Market)
	q.Set("locale", c.Cfg.Locale)

	q.Set("authorization_code", code)
	q.Set("state", state)

	u.RawQuery = q.Encode()

	return u.String()
}
