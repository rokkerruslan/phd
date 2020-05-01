package files

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
	a := App{
		assets: assets,
		opts:   opts,
		tokens: tokens.NewApp(tokens.Assets{
			assets.Db,
		}),
	}
	r := chi.NewRouter()
	r.Get("/", a.listHandler)
	r.Post("/", a.uploadHandler)
	return r
}
