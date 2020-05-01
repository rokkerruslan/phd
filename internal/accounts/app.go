package accounts

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"ph/internal/api"
	"ph/internal/tokens"
)

type (
	Resources struct {
		Db *pgxpool.Pool
	}
	Opts struct {
		GlobalSalt           []byte
		BcryptWorkFactor     int
		MinLenForNewPassword int
	}
)

type App struct {
	resources Resources
	opts      Opts

	tokens *tokens.App
}

// Setup - initialize accounts App.
func Setup(resources Resources, opts Opts) chi.Router {
	app := App{
		resources: resources,
		opts:      opts,
		tokens: tokens.NewApp(tokens.Assets{
			Db: resources.Db,
		}),
	}
	r := chi.NewRouter()
	r.Use(api.ApplicationJSON)
	r.Get("/{id}", app.retrieveHandler)
	r.Delete("/{id}", app.deleteHandler)
	r.Post("/sign-in", app.signInHandler)
	r.Post("/sign-up", app.signUpHandler)
	r.Delete("/sign-out", app.signOutHandler)
	return r
}
