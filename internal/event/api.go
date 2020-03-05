package event

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"

	"photo/internal/errors"
)

var db *pgxpool.Pool

func Mount(router chi.Router, pool *pgxpool.Pool) {
	router.Route("/api/v1/events", func(apiV1 chi.Router) {
		apiV1.Get("/", List)
		apiV1.Get("/{id}", Retrieve)
		apiV1.Post("/", Create)
		apiV1.Put("/{id}", Update)
		apiV1.Delete("/{id}", Delete)
	})

	db = pool
}

type Filter struct {
}

func NewFilterFromQuery(_ url.Values) Filter {
	return Filter{}
}

func List(w http.ResponseWriter, r *http.Request) {
	filter := NewFilterFromQuery(r.URL.Query())

	events, err := ModelList(r.Context(), filter)
	if err != nil {
		errors.APIError(w, err, http.StatusInternalServerError)
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

func Retrieve(w http.ResponseWriter, r *http.Request) {
	filter, err := newFilterRetrieve(r)
	if err != nil {
		errors.APIError(w, err, http.StatusBadRequest)
		return
	}

	event, err := modelRetrieve(r.Context(), filter)
	if err != nil {
		errors.APIError(w, err, http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(event)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Println(err)
		return
	}

	if err := validateForCreate(event); err != nil {
		errors.APIError(w, err, http.StatusBadRequest)
		return
	}

	if err := event.Create(r.Context()); err != nil {
		errors.APIError(w, err, http.StatusBadRequest)
		return
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Println(err)
		return
	}

	if err := event.ValidateForUpdate(); err != nil {
		errors.APIError(w, err, http.StatusBadRequest)
		return
	}

	if err := event.Update(r.Context()); err != nil {
		errors.APIError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	filter, err := newFilterRetrieve(r)
	if err != nil {
		errors.APIError(w, err, http.StatusBadRequest)
		return
	}

	if err := modelDelete(r.Context(), filter); err != nil {
		errors.APIError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
