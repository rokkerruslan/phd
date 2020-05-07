package files

import (
	"errors"
	"log"
	"os"

	"ph/internal/tokens"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	FilesDirName = "./images"
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
