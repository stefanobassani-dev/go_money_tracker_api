package tinkapi

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) ListCredentials(ctx context.Context, userToken string) (CredentialResponse, error) {
	path := "/api/v1/credentials/list"

	var res CredentialResponse

	err := c.call(ctx, http.MethodGet, path, "json", nil, &res, userToken)
	if err != nil {
		return CredentialResponse{}, fmt.Errorf("errore creazione list: %w", err)
	}

	return res, nil
}
