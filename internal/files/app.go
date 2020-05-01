package files

import (
	"errors"
	"log"
	"os"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"ph/internal/tokens"
)

const (
	FilesDirName = "./images"
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

	if err := os.Mkdir(FilesDirName, os.ModePerm); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal(err)
		}
	}

	r := chi.NewRouter()
	r.Get("/", a.listHandler)
	r.Post("/", a.uploadHandler)
	return r
}
