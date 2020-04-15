package offer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"ph/internal/api"
)

func (app *App) createOfferHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "App.createOfferHandler fails: %s"

	var offer Offer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	if err := offer.ValidateForCreate(); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	accountID, err := app.tokens.RetrieveAccountIDFromRequest(r.Context(), r)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	if !offer.CanBeCreated(accountID) {
		api.Error(w, fmt.Errorf(baseErr, "you can't create offer for this account"), http.StatusBadRequest)
		return
	}

	if offer, err = app.createOffer(r.Context(), offer); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(offer); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}
}

type Filter struct {
	AccountID int
	EventID   int
}

func NewFilterFromQuery(values url.Values) (Filter, error) {
	baseErr := "NewFilterFromQuery fails: %v"

	var f Filter
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

func (app *App) list(w http.ResponseWriter, r *http.Request) {
	filter, err := NewFilterFromQuery(r.URL.Query())
	if err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	offers, err := app.offerList(r.Context(), filter)
	if err != nil {
		api.Error(w, err, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(offers)
}
