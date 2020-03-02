package event

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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

func List(w http.ResponseWriter, r *http.Request) {
	events, err := ModelList(r.Context(), Filter{})
	if err != nil {
		errors.APIError(w, err, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(events)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Println(err)
		return
	}

	if err := event.Validate(); err != nil {
		errors.APIError(w, err, http.StatusBadRequest)
		return
	}

	log.Println("Create result:", event.Insert(r.Context()))
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("updated"))
}
