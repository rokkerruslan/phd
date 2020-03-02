package event

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v4"

	"photo/internal/errors"
)

// todo: use pool
var db *pgx.Conn

func init() {
	var err error
	db, err = pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:10003/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
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

func Retrieve(w http.ResponseWriter, r *http.Request) {

}

func Create(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Println(err)
		return
	}

	if err := event.ValidateForCreate(); err != nil {
		errors.APIError(w, err, http.StatusBadRequest)
		return
	}

	if err := event.Insert(r.Context()); err != nil {
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
