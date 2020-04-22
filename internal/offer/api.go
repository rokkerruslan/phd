package offer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ph/internal/api"
)

func (app *App) createHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "createHandler fails: %s"

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

	if !offer.canBeCreated(accountID) {
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

func (app *App) updateHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "updateHandler fails: %v"

	f, err := api.NewRetrieveFilter(r)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	fmt.Println("UPDATE OFFER ID:", f.ID)
}

func (app *App) list(w http.ResponseWriter, r *http.Request) {
	filter, err := NewListFilterFromQuery(r.URL.Query())
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
