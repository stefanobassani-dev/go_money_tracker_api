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

type TokenManager struct {
	clientAccessToken string
	expiry            time.Time

	httpClient *http.Client
	cfg        *config.TinkConfig
	mu         sync.RWMutex
}

func NewTokenManager(httpClient *http.Client, cfg *config.TinkConfig) *TokenManager {
	return &TokenManager{
		httpClient: httpClient,
		cfg:        cfg,
	}
}

func (m *TokenManager) GetToken(ctx context.Context) (string, error) {
	m.mu.RLock()
	if m.clientAccessToken != "" && time.Now().Before(m.expiry.Add(-1*time.Minute)) {
		m.mu.RUnlock()
		return m.clientAccessToken, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.clientAccessToken != "" && time.Now().Before(m.expiry.Add(-1*time.Minute)) {
		return m.clientAccessToken, nil
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

	newToken, expiresIn, err := m.RefreshClientToken(ctx, scopes)
	if err != nil {
		return "", err
	}

	m.clientAccessToken = newToken
	m.expiry = time.Now().Add(time.Duration(expiresIn) * time.Second)

	return m.clientAccessToken, nil
}

func (m *TokenManager) RefreshClientToken(ctx context.Context, scopes []string) (string, int, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", strings.Join(scopes, ","))
	data.Set("client_id", m.cfg.ClientId)
	data.Set("client_secret", m.cfg.ClientSecret)
	body := strings.NewReader(data.Encode())

	var res struct {
		AccessToken    string `json:"access_token"`
		TokenExpiresIn int    `json:"expires_in"`
	}

	props := Props{
		ctx:         ctx,
		httpClient:  m.httpClient,
		baseUrl:     m.cfg.BaseUrl,
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
