package event

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
	r.Get("/", a.eventListHandler)
	r.Get("/{id}", a.retrieve)
	r.Post("/", a.create)
	r.Put("/{id}", a.update)
	r.Delete("/{id}", a.delete)
	return r
}
