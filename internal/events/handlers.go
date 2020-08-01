package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"ph/internal/api"
)

func (app *App) listHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "events.listHandler fails: %v"

	f, err := api.NewAccountAndEventFilter(r.URL.Query())
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	events, err := app.eventList(r.Context(), f)
	if err != nil {
		api.Error(w, err, http.StatusInternalServerError)
		return
	}

	api.Response(w, events)
}

func (app *App) listSuggestedHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "events.listSuggestedHandler fails: %v"

	accountID, err := app.assets.Tokens.RetrieveAccountIDFromRequest(r.Context(), r)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	events, err := app.suggestedEvents(r.Context(), accountID)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	api.Response(w, events)
}

func (app *App) retrieveHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "events.retrieveHandler fails: %v"

	filter, err := api.NewRetrieveFilter(r)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	event, err := app.eventRetrieve(r.Context(), filter)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	api.Response(w, event)
}

func (app *App) createHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "events.createHandler fails: %v"

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	if err := validateForCreate(event); err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	accountID, err := app.assets.Tokens.RetrieveAccountIDFromRequest(r.Context(), r)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	if accountID != event.OwnerID {
		api.Error(w, fmt.Errorf(baseErr, "didn't authorized"), http.StatusBadRequest)
		return
	}

	event, err = app.eventCreate(r.Context(), event)
	if err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	api.Response(w, event)
}

func (app *App) updateHandler(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	var err error
	event.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	if err := event.ValidateForUpdate(); err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	if e, err := app.eventUpdate(r.Context(), event); err != nil {
		api.Error(w, err, http.StatusBadRequest)
	} else {
		api.Response(w, e)
	}
}

func (app *App) deleteHandler(w http.ResponseWriter, r *http.Request) {
	filter, err := api.NewRetrieveFilter(r)
	if err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	if err := app.eventDelete(r.Context(), filter.ID); err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
