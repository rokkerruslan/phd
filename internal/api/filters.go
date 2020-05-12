package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type RetrieveFilter struct {
	ID int
}

// TODO: Just a number? Without additional struct?
func NewRetrieveFilter(r *http.Request) (f RetrieveFilter, err error) {
	baseErr := "NewRetrieveFilter fails: %v"
	raw := chi.URLParam(r, "id")
	if raw == "" {
		return f, fmt.Errorf(baseErr, "`id` param doesn't present")
	}
	f.ID, err = strconv.Atoi(raw)
	if err != nil {
		return f, fmt.Errorf(baseErr, err)
	}
	return f, err
}
