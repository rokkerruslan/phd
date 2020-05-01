package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Response(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		Error(w, fmt.Errorf("write response fails: %v", err), http.StatusInternalServerError)
		return
	}
}
