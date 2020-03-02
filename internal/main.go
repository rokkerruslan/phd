package internal

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"photo/internal/offer"

	"photo/internal/event"
)

// Run - entry point for internal package
func Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Route("/api/v1/events", func(apiV1 chi.Router) {
		apiV1.Get("/", event.List)
		apiV1.Get("/{id}", event.Retrieve)
		apiV1.Post("/", event.Create)
		apiV1.Put("/{id}", event.Update)
	})
	r.Route("/api/v1/offers", func(apiV1 chi.Router) {
		apiV1.Get("/", offer.List)
		apiV1.Post("/", offer.Create)
	})

	log.Println("photo starts on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
