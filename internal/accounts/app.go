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

type app struct {
	resources Resources
	opts      Opts

	tokens *tokens.App
}

// Setup - initialize accounts app.
func Setup(resources Resources, opts Opts) chi.Router {
	a := app{
		resources: resources,
		opts:      opts,
		tokens: tokens.NewApp(tokens.Assets{
			Db: resources.Db,
		}),
	}
	r := chi.NewRouter()
	r.Use(api.ApplicationJSON)
	r.Get("/{id}", a.retrieveHandler)
	r.Delete("/{id}", a.deleteHandler)
	r.Post("/sign-in", a.signInHandler)
	r.Post("/sign-up", a.signUpHandler)
	r.Delete("/sign-out", a.signOutHandler)
	return r
}
