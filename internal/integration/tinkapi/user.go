package tinkapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetUserDetails(ctx context.Context, userToken string) (TinkUserResponse, error) {
	path := "/user"

	var res TinkUserResponse

	if err := c.call(ctx, http.MethodGet, path, "form", nil, &res, userToken); err != nil {
		return TinkUserResponse{}, fmt.Errorf("tinkapi client user call failed: %w", err)
	}

	return res, nil
}

func (c *Client) CreateUser(ctx context.Context, externalID string, market string, locale string, clientToken string) (CreateUserResponse, error) {
	path := "/user/create"

	reqBody := CreateUserRequest{
		ExternalUserID: externalID,
		Market:         market,
		Locale:         locale,
		RetentionClass: "permanent",
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return CreateUserResponse{}, fmt.Errorf("error marshalling user request: %w", err)
	}

	var res CreateUserResponse
	err = c.call(ctx, http.MethodPost, path, "json", bytes.NewReader(jsonData), &res, clientToken)
	if err != nil {
		return CreateUserResponse{}, err
	}

	return res, nil
}

func (c *Client) DeleteUser(ctx context.Context, userToken string) error {
	path := "/user/delete"

	if err := c.call(ctx, http.MethodPost, path, "form", nil, nil, userToken); err != nil {
		return fmt.Errorf("tinkapi client user delete call failed: %w", err)
	}

	return nil
}
