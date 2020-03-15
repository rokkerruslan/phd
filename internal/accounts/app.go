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
		GlobalSalt       []byte
		BcryptWorkFactor int
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
	r.Delete("/{id}", a.Delete)
	r.Post("/sign-in", a.SignIn)
	r.Post("/sign-up", a.SignUp)
	r.Delete("/sign-out", a.SignOut)
	return r
}
