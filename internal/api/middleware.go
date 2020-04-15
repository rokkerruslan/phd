package api

import (
	"net/http"
)

func ApplicationJSON(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
