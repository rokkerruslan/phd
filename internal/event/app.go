package event

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"ph/internal/api"
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
	r.Use(api.ApplicationJSON)
	r.Get("/", a.eventListHandler)
	r.Get("/{id}", a.retrieve)
	r.Post("/", a.create)
	r.Put("/{id}", a.update)
	r.Delete("/{id}", a.delete)
	return r
}
