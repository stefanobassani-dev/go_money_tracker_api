package tink

import (
	"context"
	"net/http"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type credential struct {
	ID                      string            `json:"id"`
	ProviderName            string            `json:"providerName"`
	Type                    string            `json:"type"`
	Status                  string            `json:"status"`
	StatusPayload           string            `json:"statusPayload"`
	StatusUpdated           int64             `json:"statusUpdated"`
	Updated                 int64             `json:"updated"`
	SessionExpiryDate       int64             `json:"sessionExpiryDate"`
	UserID                  string            `json:"userId"`
	Fields                  map[string]string `json:"fields"`
	SupplementalInformation interface{}       `json:"supplementalInformation"`
}

func (c *Client) GetUserCredential(ctx context.Context, credentialID string,
	externalUserID string) (domain.Credential, error) {
	scopes := []string{"credentials:write", "credentials:read"}
	accessToken, err := c.ExchangeUserToken(ctx, externalUserID, scopes)
	if err != nil {
		return domain.Credential{}, err
	}

	var res credential
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
