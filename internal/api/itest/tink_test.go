package itest

import (
	"net/http"
	"testing"
)

func TestGetConnectURL(t *testing.T) {
	tests := []Test{
		{
			name:           "Error con jwt mancante",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	runTests(t, "/tink/link", http.MethodGet, tests)
}
