package api

import (
	"encoding/json"
	"fmt"
	"net/http"

)

func Response(w http.ResponseWriter, obj interface{}) {
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		Error(w, fmt.Errorf("write response fails: %v", err), http.StatusInternalServerError)
		return
	}
}
