package offer

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	Assets struct {
		Db *pgxpool.Pool
	}
	Opts struct{}
)

type App struct {
	assets Assets
	opts   Opts
}

func Setup(assets Assets, opts Opts) chi.Router {
	a := App{
		assets: assets,
		opts:   opts,
	}
	r := chi.NewRouter()
	r.Get("/", a.list)
	r.Post("/", a.createOfferHandler)
	return r
}
