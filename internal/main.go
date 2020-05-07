package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"ph/internal/accounts"
	"ph/internal/event"
	"ph/internal/files"
	"ph/internal/offer"
	"ph/internal/tokens"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Run - entry point for internal package
func Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders: []string{"X-Auth-Token"},
		Debug:          true,
	}))

	opts, err := newOptions()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, opts.databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	tokenApp := tokens.NewApp(
		tokens.Assets{
			Db: pool,
		},
		tokens.Opts{
			TokenTTL: opts.tokenTTL,
		},
	)

	// We mount all our sub-applications for root
	// router. Consistency isn't important.
	r.Mount("/events", event.Setup(
		event.Assets{
			Db:     pool,
			Tokens: tokenApp,
		},
		event.Opts{},
	))
	r.Mount("/offers", offer.Setup(
		offer.Assets{
			Db:     pool,
			Tokens: tokenApp,
		},
		offer.Opts{},
	))
	r.Mount("/accounts", accounts.Setup(
		accounts.Assets{
			Db:     pool,
			Tokens: tokenApp,
		},
		accounts.Opts{
			GlobalSalt:           opts.globalSalt,
			BcryptWorkFactor:     opts.bcryptWorkFactor,
			MinLenForNewPassword: opts.minLenForNewPassword,
		},
	))
	r.Mount("/files", files.Setup(
		files.Assets{
			Db:     pool,
			Tokens: tokenApp,
		},
		files.Opts{},
	))

	log.Println(fmt.Sprintf("daemon bind socket on %s", opts.addr))
	if err := http.ListenAndServe(opts.addr, r); err != nil {
		log.Fatal(err)
	}
}
