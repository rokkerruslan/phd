package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"ph/internal/accounts"
	"ph/internal/event"
	"ph/internal/files"
	"ph/internal/offer"
)

// Run - entry point for internal package
func Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
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

	// We mount all our sub-applications for root
	// router. Consistency isn't important.
	r.Mount("/events", event.Setup(
		event.Resources{
			Db: pool,
		},
		event.Opts{},
	))
	r.Mount("/offers", offer.Setup(
		offer.Assets{
			Db: pool,
		},
		offer.Opts{},
	))
	r.Mount("/accounts", accounts.Setup(
		accounts.Resources{
			Db: pool,
		},
		accounts.Opts{
			GlobalSalt:           opts.globalSalt,
			BcryptWorkFactor:     opts.bcryptWorkFactor,
			MinLenForNewPassword: opts.minLenForNewPassword,
		},
	))
	r.Mount("/files", files.Setup(
		files.Assets{
			Db: pool,
		},
		files.Opts{},
	))

	log.Println(fmt.Sprintf("daemon bind socket on %s", opts.addr))
	if err := http.ListenAndServe(opts.addr, r); err != nil {
		log.Fatal(err)
	}
}
