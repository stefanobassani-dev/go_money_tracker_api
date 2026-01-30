package tink

import (
	"fmt"
	"net/http"

	"github.com/stefanobassani-dev/money-tracker/internal/config"
)

type APIPath string

const (
	PathTokenExchange     APIPath = "/api/v1/oauth/token"
	PathUserCreate        APIPath = "/api/v1/user/create"
	PathGetUser           APIPath = "/api/v1/user"
	PathAuthorize         APIPath = "/api/v1/oauth/authorization-grant"
	PathAuthorizeDelegate APIPath = "/api/v1/oauth/authorization-grant/delegate"
	PathGetCredential     APIPath = "/api/v1/credentials/"
	PathListAccounts      APIPath = "/data/v2/accounts"
	PathProviderConsent   APIPath = "/api/v1/provider-consents"
)

type ContentType string

const (
	ContentTypeJSON ContentType = "application/json"
	ContentTypeForm ContentType = "application/x-www-form-urlencoded"
)

const TinkActorClientID = "df05e4b379934cd09963197cc855bfe9"
const ApiConnectURL = "https://link.tink.com/1.0/transactions/connect-accounts"

type Client struct {
	Cfg          *config.TinkConfig
	TokenManager *TokenManager
	httpClient   *http.Client
}

func NewTinkClient(cfg *config.TinkConfig, tokenManager *TokenManager, http *http.Client) *Client {
	return &Client{
		Cfg:          cfg,
		TokenManager: tokenManager,
		httpClient:   http,
	}
}

type TinkError struct {
	StatusCode int
	Code       string
	Message    string
	TrackingID string
}

func (e *TinkError) Error() string {
	return fmt.Sprintf("tinkapi api error: %s (status: %d, tracking: %s)", e.Message, e.StatusCode, e.TrackingID)
}
