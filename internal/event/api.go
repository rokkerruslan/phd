package event

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi"

	"ph/internal/api"
)

type Filter struct {
}

func NewFilterFromQuery(_ url.Values) Filter {
	return Filter{}
}

func (app *app) listHandler(w http.ResponseWriter, r *http.Request) {
	filter := NewFilterFromQuery(r.URL.Query())

	events, err := app.eventList(r.Context(), filter)
	if err != nil {
		api.Error(w, err, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(events)
}

func (app *app) retrieveHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *app) createHandler(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	if err := validateForCreate(event); err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	if err := app.CreateEvent(r.Context(), event); err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}
}

func (app *app) updateHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *app) deleteHandler(w http.ResponseWriter, r *http.Request) {
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
