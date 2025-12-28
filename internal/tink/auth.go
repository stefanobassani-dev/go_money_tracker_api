package tink

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const TinkLinkActorClientID = "df05e4b379934cd09963197cc855bfe9"

func (c *Client) AuthorizationGrant(externalID string, clientToken string) (string, error) {
	data := url.Values{}
	data.Set("external_user_id", externalID)
	data.Set("scope", "user:read,user:delete")

	var res struct {
		Code string `json:"code"`
	}

	err := c.call("POST", "/oauth/authorization-grant", "form", strings.NewReader(data.Encode()), &res, clientToken)
	if err != nil {
		return "", err
	}
	return res.Code, nil
}

func (c *Client) AuthorizationGrantDelegate(externalUserID string, clientToken string) (string, error) {
	path := "/oauth/authorization-grant/delegate"

	data := url.Values{}
	data.Set("actor_client_id", TinkLinkActorClientID)
	data.Set("external_user_id", externalUserID)
	//TODO change maybe with email
	data.Set("id_hint", externalUserID)

	scopes := []string{
		"authorization:read",
		"authorization:grant",
		"credentials:refresh",
		"credentials:read",
		"credentials:write",
		"providers:read",
		"user:read",
		"accounts:read",
		"transactions:read",
	}
	data.Set("scope", strings.Join(scopes, ","))

	var res struct {
		Code string `json:"code"`
	}

	err := c.call("POST", path, "form", strings.NewReader(data.Encode()), &res, clientToken)
	if err != nil {
		return "", fmt.Errorf("delegate error: %w", err)
	}

	return res.Code, nil
}

func (c *Client) GetUserAccessToken(code string) (*TokenResponse, error) {
	path := "/oauth/token"

	data := url.Values{}
	data.Set("client_id", c.Cfg.ClientId)
	data.Set("client_secret", c.Cfg.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)

	var res TokenResponse
	err := c.call("POST", path, "form", strings.NewReader(data.Encode()), &res, "")
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetClientAccessToken() (string, error) {
	path := "/oauth/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "user:create,authorization:grant,user:read,user:delete")
	data.Set("client_id", c.Cfg.ClientId)
	data.Set("client_secret", c.Cfg.ClientSecret)
	body := strings.NewReader(data.Encode())

	var res struct {
		UserAccessToken string `json:"access_token"`
	}

	if err := c.call(http.MethodPost, path, "form", body, &res, ""); err != nil {
		return "", err
	}

	return res.UserAccessToken, nil
}

func (c *Client) GetAuthorizationURL(externalUserId string, clientToken string,
	market string, locale string) (string, error) {
	code, err := c.AuthorizationGrantDelegate(externalUserId, clientToken)
	if err != nil {
		return "", err
	}

	//TODO HMAC per verificare integrità e non manipolazione
	state := externalUserId
	return c.BuildUrl(c.Cfg.ClientId, state, c.Cfg.RedirectUri, code, market, locale), nil
}

func (c *Client) BuildUrl(clientID string, state string, redirectURI string,
	code string, market string, locale string) string {
	return "https://link.tink.com/1.0/transactions/connect-accounts?client_id=" + clientID +
		"&state=" + state + "&redirect_uri=" + redirectURI + "&authorization_code=" + code +
		"&market=" + market + "&locale=" + locale
}
