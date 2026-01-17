package tink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Props struct {
	ctx         context.Context
	httpClient  *http.Client
	baseUrl     string
	method      string
	path        string
	contentType ContentType
	body        io.Reader
	result      any
	token       string
}

func call(p Props) error {
	fullURL := strings.TrimSuffix(p.baseUrl, "/") + "/" + strings.TrimPrefix(string(p.path), "/")

	req, err := http.NewRequestWithContext(p.ctx, p.method, fullURL, p.body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	if p.token != "" {
		req.Header.Set("Authorization", "Bearer "+p.token)
	}

	if p.contentType == ContentTypeJSON {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		tErr := &TinkError{
			StatusCode: resp.StatusCode,
			TrackingID: resp.Header.Get("X-Tink-Tracking-Id"),
		}

		var errorPayload struct {
			ErrorMessage string `json:"errorMessage"`
			ErrorCode    string `json:"errorCode"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errorPayload); err == nil {
			tErr.Message = errorPayload.ErrorMessage
			tErr.Code = errorPayload.ErrorCode
		} else {
			tErr.Message = "Unknown Tink API error"
		}

		return tErr
	}

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("tinkapi api error [status %d]: %s", resp.StatusCode, string(bodyBytes))
	}

	if p.result != nil {
		if err := json.NewDecoder(resp.Body).Decode(p.result); err != nil {
			return fmt.Errorf("errore decodifica risposta: %w", err)
		}
	}

	return nil
}

func toJSONReader(v any) (io.Reader, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}
