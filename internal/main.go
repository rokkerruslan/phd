package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"ph/internal/accounts"
	"ph/internal/events"
	"ph/internal/files"
	"ph/internal/offers"
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

	opts, err := FromFile()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("OPTS: %+v\n", opts)

	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, opts.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	tokenApp := tokens.NewApp(
		tokens.Assets{
			Db: pool,
		},
		tokens.Opts{
			TokenTTL: opts.TokenTTL,
		},
	)

	// We mount all our sub-applications for root
	// router. Consistency isn't important.
	r.Mount("/events", events.Setup(
		events.Assets{
			Db:     pool,
			Tokens: tokenApp,
		},
		events.Opts{},
	))
	r.Mount("/offers", offers.Setup(
		offers.Assets{
			Db:     pool,
			Tokens: tokenApp,
		},
		offers.Opts{},
	))
	r.Mount("/accounts", accounts.Setup(
		accounts.Assets{
			Db:     pool,
			Tokens: tokenApp,
		},
		accounts.Opts{
			GlobalSalt:           []byte(opts.GlobalSalt),
			BcryptWorkFactor:     opts.BcryptWorkFactor,
			MinLenForNewPassword: opts.MinLenForNewPassword,
		},
	))
	r.Mount("/files", files.Setup(
		files.Assets{
			Db:     pool,
			Tokens: tokenApp,
		},
		files.Opts{},
	))

	log.Println(fmt.Sprintf("daemon bind socket on %s", opts.Addr))
	if err := http.ListenAndServe(opts.Addr, r); err != nil {
		log.Fatal(err)
	}
}
