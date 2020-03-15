package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"photo/internal/accounts"
	"photo/internal/event"
	"photo/internal/offer"
)

// Run - entry point for internal package
func Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

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
	r.Route("/api/v1", func(apiV1 chi.Router) {
		apiV1.Mount("/events", event.Setup(
			event.Resources{
				Db: pool,
			},
			event.Options{},
		))
		apiV1.Mount("/offers", offer.Setup(
			offer.Resources{
				Db: pool,
			},
			offer.Options{},
		))
		apiV1.Mount("/accounts", accounts.Setup(
			accounts.Resources{
				Db: pool,
			},
			accounts.Options{
				GlobalSalt:       opts.globalSalt,
				BcryptWorkFactor: opts.bcryptWorkFactor,
			},
		))
	})

	log.Println(fmt.Sprintf("daemon bind socket on %s", opts.addr))
	if err := http.ListenAndServe(opts.addr, r); err != nil {
		log.Fatal(err)
	}
}
