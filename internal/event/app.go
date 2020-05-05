package event

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"ph/internal/tokens"
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

	tokens *tokens.App
}

func Setup(assets Assets, opts Opts) chi.Router {
	app := App{
		assets: assets,
		opts:   opts,
		tokens: tokens.NewApp(tokens.Assets{
			Db: assets.Db,
		}),
	}
	r := chi.NewRouter()
	r.Get("/", app.listHandler)
	r.Get("/{id}", app.retrieveHandler)
	r.Post("/", app.createHandler)
	r.Put("/{id}", app.updateHandler)
	r.Delete("/{id}", app.deleteHandler)
	return r
}
