package accounts

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	Resources struct {
		Db *pgxpool.Pool
	}
	Options struct {
		GlobalSalt           []byte
		BcryptWorkFactor     int
		MinLenForNewPassword int
	}
)

type app struct {
	resources Resources
	options   Options
}

func Setup(resources Resources, options Options) chi.Router {
	a := app{
		resources: resources,
		options:   options,
	}
	r := chi.NewRouter()
	r.Get("/{id}", a.retrieveHandler)
	r.Delete("/{id}", a.deleteHandler)
	r.Post("/sign-in", a.signInHandler)
	r.Post("/sign-up", a.signUpHandler)
	r.Delete("/sign-out", a.signOutHandler)
	return r
}
