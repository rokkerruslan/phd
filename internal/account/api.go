package account

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Mount(r chi.Router, pool *pgxpool.Pool) {
	r.Get("/{id}", Retrieve)
	r.Delete("/{id}", Delete)
}

func Retrieve(w http.ResponseWriter, r *http.Request) {

}

func Delete(w http.ResponseWriter, r *http.Request) {

}
