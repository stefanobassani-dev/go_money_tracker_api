package tink

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetClientToken(ctx context.Context, scopes []string) (string, int, error) {
	path := "/api/v1/oauth/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", strings.Join(scopes, ","))
	data.Set("client_id", c.Cfg.ClientId)
	data.Set("client_secret", c.Cfg.ClientSecret)
	body := strings.NewReader(data.Encode())

	var res struct {
		AccessToken    string `json:"access_token"`
		TokenExpiresIn int    `json:"expires_in"`
	}

	if err := c.call(ctx, http.MethodPost, path, "form", body, &res, ""); err != nil {
		return "", 0, err
	}

	return res.AccessToken, res.TokenExpiresIn, nil
}
