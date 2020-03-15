package offer

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
	r.Get("/", a.list)
	r.Post("/", a.create)
	return r
}
