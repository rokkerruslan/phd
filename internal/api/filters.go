package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

type AccountAndEventFilter struct {
	AccountID int
	EventID   int
}

func NewAccountAndEventFilter(values url.Values) (AccountAndEventFilter, error) {
	baseErr := "NewAccountAndEventFilter fails: %v"

	var f AccountAndEventFilter
	var err error
	var errors []string

	accountIDRaw := values.Get("account_id")
	if accountIDRaw != "" {
		if f.AccountID, err = strconv.Atoi(accountIDRaw); err != nil {
			errors = append(errors, fmt.Sprintf("account_id parsing fails: %v", err))
		}
	}

	eventIDRaw := values.Get("event_id")
	if eventIDRaw != "" {
		if f.EventID, err = strconv.Atoi(eventIDRaw); err != nil {
			errors = append(errors, fmt.Sprintf("event_id parsing fails: %v", err))
		}
	}

	if len(errors) != 0 {
		return f, fmt.Errorf(baseErr, strings.Join(errors, ", "))
	}

	return f, nil
}
