package api

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string
}

// TODO: logging here?
func Error(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}
