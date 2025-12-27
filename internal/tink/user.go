package tink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetUserDetails(userToken string) (TinkUserResponse, error) {
	path := "/user"

	var res TinkUserResponse

	if err := c.call(http.MethodGet, path, "form", nil, &res, userToken); err != nil {
		return TinkUserResponse{}, fmt.Errorf("tink client user call failed: %w", err)
	}

	return res, nil
}

func (c *Client) CreateUser(externalID string, market string, locale string, clientToken string) (CreateUserResponse, error) {
	path := "/user/create"

	reqBody := CreateUserRequest{
		ExternalUserID: externalID,
		Market:         market,
		Locale:         locale,
		RetentionClass: "permanent",
	}
	jsonData, marshalErr := json.Marshal(reqBody)
	if marshalErr != nil {
		return CreateUserResponse{}, fmt.Errorf("error marshalling user request: %w", marshalErr)
	}

	var res CreateUserResponse
	callErr := c.call(http.MethodPost, path, "json", bytes.NewReader(jsonData), &res, clientToken)
	if callErr != nil {
		return CreateUserResponse{}, fmt.Errorf("create user error: %w", callErr)
	}

	return res, nil
}

func (c *Client) DeleteUser(userToken string) error {
	path := "/user/delete"

	if err := c.call(http.MethodPost, path, "form", nil, nil, userToken); err != nil {
		return fmt.Errorf("tink client user delete call failed: %w", err)
	}

	return nil
}
