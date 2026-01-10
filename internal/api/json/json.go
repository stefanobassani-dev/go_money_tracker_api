package json

import (
	"encoding/json"
	"net/http"
)

func Decode(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func Success(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func Error(w http.ResponseWriter, status int, message string) {
	Success(w, status, map[string]string{"error": message})
}
