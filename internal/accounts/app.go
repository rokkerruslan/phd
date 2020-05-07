package accounts

import (
	"ph/internal/tokens"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	Assets struct {
		Db *pgxpool.Pool
	}
	Opts struct {
		GlobalSalt           []byte
		BcryptWorkFactor     int
		MinLenForNewPassword int
	}
)

type App struct {
	assets Assets
	opts   Opts

	tokens *tokens.App
}

// Setup - initialize accounts App.
func Setup(assets Assets, opts Opts) chi.Router {
	app := App{
		assets: assets,
		opts:   opts,
		tokens: tokens.NewApp(tokens.Assets{
			Db: assets.Db,
		}),
	}
	r := chi.NewRouter()
	r.Get("/{id}", app.retrieveHandler)
	r.Delete("/{id}", app.deleteHandler)
	r.Post("/sign-in", app.signInHandler)
	r.Post("/sign-up", app.signUpHandler)
	r.Delete("/sign-out", app.signOutHandler)
	return r
}
