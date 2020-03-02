package errors

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string
}

func APIError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}
