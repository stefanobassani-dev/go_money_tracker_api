package tink

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient   http.Client
	ClientId     string
	clientSecret string
	baseUrl      string
}

func NewTinkClient(clientId string, clientSecret string, baseUrl string) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: time.Second * 5,
		},
		ClientId:     clientId,
		clientSecret: clientSecret,
		baseUrl:      baseUrl,
	}
}

func (c *Client) call(method, path string, contentType string, body io.Reader, result any, token string) error {
	fullURL := strings.TrimSuffix(c.baseUrl, "/") + "/" + strings.TrimPrefix(path, "/")

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return fmt.Errorf("errore creazione richiesta: %w", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	if contentType == "json" {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("errore invio richiesta: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("tink api error [status %d]: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("errore decodifica risposta: %w", err)
		}
	}

	return nil
}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
