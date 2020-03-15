package offer

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	Resources struct {
		Db *pgxpool.Pool
	}
	Options struct{}
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
	r.Get("/", a.list)
	r.Post("/", a.create)
	return r
}
