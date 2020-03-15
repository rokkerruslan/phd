package offer

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Resources struct {
	Db *pgxpool.Pool
}

type Options struct{}

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
	r.Get("/", app.list)
	r.Post("/", app.create)
	return r
}
