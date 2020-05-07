package offer

import (
	"ph/internal/tokens"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	Assets struct {
		Db     *pgxpool.Pool
		Tokens *tokens.App
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
	r.Post("/", a.createHandler)
	r.Put("/{id}", a.updateHandler)
	r.Delete("/{id}", a.deleteHandler)
	return r
}
