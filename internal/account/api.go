package account

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Options struct {
	GlobalSalt []byte
}

var db *pgxpool.Pool
var options Options

func Mount(r chi.Router, pool *pgxpool.Pool, o Options) {
	r.Get("/{id}", Retrieve)
	r.Delete("/{id}", Delete)

	db = pool
	options = o
}

func Retrieve(w http.ResponseWriter, r *http.Request) {

}

func Delete(w http.ResponseWriter, r *http.Request) {

}
