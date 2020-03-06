package auth

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Mount(router chi.Router, pool *pgxpool.Pool) {
	router.Route("/api/v1/auth", func(apiV1 chi.Router) {
		apiV1.Post("/sign-in", SignIn)
		apiV1.Post("/sign-up", SignUp)
		apiV1.Delete("/sign-out", SignOut)
	})
}
