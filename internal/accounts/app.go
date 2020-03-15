package accounts

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Resources struct {
	Db *pgxpool.Pool
}

type Options struct {
	GlobalSalt       []byte
	BcryptWorkFactor int
}

type App struct {
	resources Resources
	options   Options
}

func Setup(resources Resources, options Options) chi.Router {
	app := App{
		resources: resources,
		options:   options,
	}
	r := chi.NewRouter()
	r.Get("/{id}", app.Retrieve)
	r.Delete("/{id}", app.Delete)

	r.Post("/sign-in", app.SignIn)
	r.Post("/sign-up", app.SignUp)
	r.Delete("/sign-out", app.SignOut)
	return r
}
