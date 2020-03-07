package auth

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Options struct {
	GlobalSalt       []byte
	BcryptWorkFactor int
}

var db *pgxpool.Pool
var options Options

func Mount(r chi.Router, pool *pgxpool.Pool, opts Options) {
	r.Post("/sign-in", SignIn)
	r.Post("/sign-up", SignUp)
	r.Delete("/sign-out", SignOut)

	db = pool
	options = opts
}
