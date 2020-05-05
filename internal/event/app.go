package event

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	Resources struct {
		Db *pgxpool.Pool
	}
	Opts struct{}
)

type app struct {
	resources Resources
	opts      Opts
}

func Setup(resources Resources, opts Opts) chi.Router {
	a := app{
		resources: resources,
		opts:      opts,
	}
	r := chi.NewRouter()
	r.Get("/", a.listHandler)
	r.Get("/{id}", a.retrieveHandler)
	r.Post("/", a.createHandler)
	r.Put("/{id}", a.updateHandler)
	r.Delete("/{id}", a.deleteHandler)
	return r
}
