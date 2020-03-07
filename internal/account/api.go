package account

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"photo/internal/errors"
	"photo/internal/session"
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

const AuthTokenName = "X-Auth-Token"

func Retrieve(w http.ResponseWriter, r *http.Request) {
	baseErr := "account.Retrieve fails: %v"

	token := r.Header.Get(AuthTokenName)
	if token == "" {
		errors.APIError(w, fmt.Errorf(baseErr, fmt.Sprintf("`%s` isn't set", AuthTokenName)), http.StatusForbidden)
		return
	}

	id, err := session.RetrieveSession(r.Context(), db, token)
	if err != nil {
		errors.APIError(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	acc, err := RetrieveByID(r.Context(), id)
	if err != nil {
		errors.APIError(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(acc)
}

func Delete(_ http.ResponseWriter, _ *http.Request) {

}
