package event

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
	app := App{
		assets: assets,
		opts:   opts,
	}
	r := chi.NewRouter()
	r.Get("/", app.listHandler)
	r.Get("/suggested", app.listSuggestedHandler)
	r.Get("/{id}", app.retrieveHandler)
	r.Post("/", app.createHandler)
	r.Put("/{id}", app.updateHandler)
	r.Delete("/{id}", app.deleteHandler)
	return r
}
