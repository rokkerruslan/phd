package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"photo/internal/account"
	"photo/internal/event"
	"photo/internal/offer"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Run - entry point for internal package
func Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	opts, err := newOptions()
	if err != nil {
		log.Fatal(err)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, opts.databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	// We mount all our sub-applications for root
	// router. Consistency isn't important.
	r.Route("/api/v1", func(apiV1 chi.Router) {
		apiV1.Route("/events", func(r chi.Router) {
			event.Mount(r, pool)
		})
		apiV1.Route("/offers", func(r chi.Router) {
			offer.Mount(r, pool)
		})
		apiV1.Mount("/accounts", account.NewApp(
			account.Resources{
				Db: pool,
			},
			account.Options{
				GlobalSalt:       opts.globalSalt,
				BcryptWorkFactor: opts.bcryptWorkFactor,
			}),
		)
	})

	log.Println(fmt.Sprintf("daemon bind socket on %s", opts.addr))
	if err := http.ListenAndServe(opts.addr, r); err != nil {
		log.Fatal(err)
	}
}
