package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"sprayer/internal/api"
	"sprayer/internal/job"
	"sprayer/internal/profile"
)

func main() {
	godotenv.Load()

	store, err := job.NewStore()
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	pStore, err := profile.NewStore(store.DB)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	h := api.NewHandler(store, pStore)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", h.HealthCheck)
		
		r.Route("/jobs", func(r chi.Router) {
			r.Get("/", h.ListJobs)
			r.Post("/scrape", h.ScrapeJobs)
		})

		r.Route("/profiles", func(r chi.Router) {
			r.Get("/", h.ListProfiles)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
