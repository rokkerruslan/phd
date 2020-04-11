package offer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"ph/internal/api"
	"ph/internal/tokens"
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

	accountID, err := tokens.RetrieveAccountIDByToken(r.Context(), app.assets.Db, r)
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
}

func NewFilterFromQuery(values url.Values) (f Filter, err error) {
	defErr := "filter creating fails: %v"
	if f.AccountID, err = strconv.Atoi(values.Get("account_id")); err != nil {
		return f, fmt.Errorf(defErr, fmt.Errorf("account_id parsing fails: %v", err))
	}
	return f, nil
}

func (app *App) list(w http.ResponseWriter, r *http.Request) {
	filter, err := NewFilterFromQuery(r.URL.Query())
	if err != nil {
		api.Error(w, fmt.Errorf("account_id parsing fails: %v", err), http.StatusBadRequest)
		return
	}

	offers, err := app.offerList(r.Context(), filter)
	if err != nil {
		api.Error(w, err, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(offers)
}
