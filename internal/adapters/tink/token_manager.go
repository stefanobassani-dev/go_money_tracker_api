package tink

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/stefanobassani-dev/money-tracker/internal/config"
)

type ClientTokenManager interface {
	GetToken(ctx context.Context) (string, error)
}

// simple client token manager with one token with all scopes (highly insecure)
type SimpleTokenManager struct {
	clientAccessToken string
	expiry            time.Time

	httpClient *http.Client
	cfg        *config.TinkConfig
	mu         sync.RWMutex
}

func NewSimpleTokenManager(httpClient *http.Client, cfg *config.TinkConfig) *SimpleTokenManager {
	return &SimpleTokenManager{
		httpClient: httpClient,
		cfg:        cfg,
	}
}

func (sm *SimpleTokenManager) GetToken(ctx context.Context) (string, error) {
	sm.mu.RLock()
	if sm.clientAccessToken != "" && time.Now().Before(sm.expiry.Add(-1*time.Minute)) {
		sm.mu.RUnlock()
		return sm.clientAccessToken, nil
	}
	sm.mu.RUnlock()

	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.clientAccessToken != "" && time.Now().Before(sm.expiry.Add(-1*time.Minute)) {
		return sm.clientAccessToken, nil
	}

	scopes := []string{
		"authorization:read",
		"authorization:grant",
		"credentials:refresh",
		"credentials:read",
		"credentials:write",
		"providers:read",
		"user:read",
		"user:create",
		"accounts:read",
		"transactions:read"}

	newToken, expiresIn, err := sm.RefreshClientToken(ctx, scopes)
	if err != nil {
		return "", err
	}

	sm.clientAccessToken = newToken
	sm.expiry = time.Now().Add(time.Duration(expiresIn) * time.Second)

	return sm.clientAccessToken, nil
}

func (sm *SimpleTokenManager) RefreshClientToken(ctx context.Context, scopes []string) (string, int, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", strings.Join(scopes, ","))
	data.Set("client_id", sm.cfg.ClientId)
	data.Set("client_secret", sm.cfg.ClientSecret)
	body := strings.NewReader(data.Encode())

	var res struct {
		AccessToken    string `json:"access_token"`
		TokenExpiresIn int    `json:"expires_in"`
	}

	props := Props{
		ctx:         ctx,
		httpClient:  sm.httpClient,
		baseUrl:     sm.cfg.BaseUrl,
		path:        string(PathTokenExchange),
		method:      http.MethodPost,
		body:        body,
		contentType: "form",
		result:      &res,
		token:       "",
	}

	if err := call(props); err != nil {
		return "", 0, err
	}

	return res.AccessToken, res.TokenExpiresIn, nil
}
