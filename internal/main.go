package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"photo/internal/account"
	"photo/internal/auth"
	"photo/internal/event"
	"photo/internal/offer"

	"github.com/jackc/pgx/v4/pgxpool"
)

type app struct {
	pool *pgxpool.Pool
	r    chi.Router
}

// Run - entry point for internal package
func Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	opts, err := newOptions()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("start with options: %+v", opts)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, opts.databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	event.Mount(r, pool)

	r.Route("/api/v1/offers", func(apiV1 chi.Router) {
		apiV1.Get("/", offer.List)
		apiV1.Post("/", offer.Create)
	})
	r.Route("/api/v1/accounts", func(apiV1 chi.Router) {
		apiV1.Get("/{id}", account.Retrieve)
		apiV1.Delete("/{id}", account.Delete)
	})
	r.Route("/api/v1/auth", func(apiV1 chi.Router) {
		apiV1.Post("/sign-in", auth.SignIn)
		apiV1.Post("/sign-up", auth.SignUp)
	})

	log.Println(fmt.Sprintf("daemon bind socket on %s", opts.addr))
	if err := http.ListenAndServe(opts.addr, r); err != nil {
		log.Fatal(err)
	}
}
