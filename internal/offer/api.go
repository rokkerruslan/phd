package offer

import (
	"encoding/json"
	"errors"
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

	var eventOwnerID int
	if err := app.assets.Db.
		QueryRow(r.Context(), "SELECT owner_id FROM events WHERE id = $1", offer.EventID).
		Scan(&eventOwnerID); err != nil {
			api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
			return
	}

	if eventOwnerID == accountID {
		api.Error(w, fmt.Errorf(baseErr, "you can't create offer for you event"), http.StatusBadRequest)
		return
	}

	if offer, err = app.createOffer(r.Context(), offer); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	api.Response(w, offer)
}

func (app *App) updateHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "updateHandler fails: %v"

	f, err := api.NewRetrieveFilter(r)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	var eventOwnerID int
	if err := app.assets.Db.QueryRow(
		r.Context(),
		"SELECT owner_id FROM events JOIN offers ON events.id = offers.event_id WHERE offers.id = $1",
		f.ID,
	).Scan(&eventOwnerID); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	accountID, err := app.tokens.RetrieveAccountIDFromRequest(r.Context(), r)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	if accountID != eventOwnerID {
		api.Error(w, fmt.Errorf(baseErr, errors.New("only Event Owner can update order")), http.StatusBadRequest)
		return
	}

	var offer Offer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	offer.ID = f.ID
	offer, err = app.updateOffer(r.Context(), offer)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	api.Response(w, offer)
}

func (app *App) deleteHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "updateHandler fails: %v"

	f, err := api.NewRetrieveFilter(r)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	accountID, err := app.tokens.RetrieveAccountIDFromRequest(r.Context(), r)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	var offerOwnerID int
	if err := app.assets.Db.QueryRow(
		r.Context(),
		"SELECT account_id FROM offers WHERE id = $1",
		f.ID,
	).Scan(&offerOwnerID); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	if accountID != offerOwnerID {
		api.Error(w, fmt.Errorf(baseErr, "only owner can delete offer"), http.StatusBadRequest)
		return
	}

	if err := app.deleteOffer(r.Context(), f); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

	api.Response(w, offers)
}
