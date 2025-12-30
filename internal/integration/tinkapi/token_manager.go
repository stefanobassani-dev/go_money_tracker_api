package tinkapi

import (
	"context"
	"sync"
	"time"
)

type TokenManager struct {
	client *Client

	clientAccessToken string
	expiry            time.Time

	mu sync.RWMutex
}

func NewTokenManager(c *Client) *TokenManager {
	return &TokenManager{
		client: c,
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

	newToken, expiresIn, err := m.client.GetClientAccessToken(ctx)
	if err != nil {
		return "", err
	}

	m.clientAccessToken = newToken
	m.expiry = time.Now().Add(time.Duration(expiresIn) * time.Second)

	return m.clientAccessToken, nil
}
