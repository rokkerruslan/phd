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

type Filter struct {
}

func NewFilterFromQuery(_ url.Values) Filter {
	return Filter{}
}

func (app *app) eventListHandler(w http.ResponseWriter, r *http.Request) {
	filter := NewFilterFromQuery(r.URL.Query())

	events, err := app.eventList(r.Context(), filter)
	if err != nil {
		api.Error(w, err, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(events)
}

type filterRetrieve struct {
	ID int
}

func newFilterRetrieve(r *http.Request) (f filterRetrieve, err error) {
	baseErr := "newFilterRetrieve fails: %v"
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

func (app *app) retrieve(w http.ResponseWriter, r *http.Request) {
	filter, err := newFilterRetrieve(r)
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

func (app *app) create(w http.ResponseWriter, r *http.Request) {
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

func (app *app) update(w http.ResponseWriter, r *http.Request) {
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

func (app *app) delete(w http.ResponseWriter, r *http.Request) {
	filter, err := newFilterRetrieve(r)
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
