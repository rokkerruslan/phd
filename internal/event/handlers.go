package event

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi"

	"ph/internal/api"
)

type filter struct {
}

func newFilterFromQuery(_ url.Values) filter {
	return filter{}
}

func (app *App) listHandler(w http.ResponseWriter, r *http.Request) {
	filter := newFilterFromQuery(r.URL.Query())

	events, err := app.eventList(r.Context(), filter)
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

	api.Response(w, struct{ Ok int }{Ok: accountID})
}

func (app *App) retrieveHandler(w http.ResponseWriter, r *http.Request) {
	filter, err := api.NewRetrieveFilter(r)
	if err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	event, err := app.retrieveEvent(r.Context(), filter)
	if err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(event)
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

	event, err = app.createEvent(r.Context(), event)
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

	if err := app.updateEvent(r.Context(), event); err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *App) deleteHandler(w http.ResponseWriter, r *http.Request) {
	filter, err := api.NewRetrieveFilter(r)
	if err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	if err := app.deleteEvent(r.Context(), filter); err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
