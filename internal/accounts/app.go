package accounts

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	Resources struct {
		Db *pgxpool.Pool
	}
	Opts struct {
		GlobalSalt           []byte
		BcryptWorkFactor     int
		MinLenForNewPassword int
	}
)

type app struct {
	resources Resources
	opts      Opts
}

// Setup - initialize accounts app.
func Setup(resources Resources, opts Opts) chi.Router {
	a := app{
		resources: resources,
		opts:      opts,
	}
	r := chi.NewRouter()
	r.Get("/{id}", a.retrieveHandler)
	r.Delete("/{id}", a.deleteHandler)
	r.Post("/sign-in", a.signInHandler)
	r.Post("/sign-up", a.signUpHandler)
	r.Delete("/sign-out", a.signOutHandler)
	return r
}
