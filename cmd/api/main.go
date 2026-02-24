package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"sprayer/src/api"
	"sprayer/src/api/job"
	"sprayer/src/api/profile"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	if envPort := os.Getenv("PORT"); envPort != "" {
		*port = envPort
	}

	// Initialize stores
	jobStore, err := job.NewStore()
	if err != nil {
		log.Fatalf("Failed to initialize job store: %v", err)
	}
	profileStore, err := profile.NewStore(jobStore.DB)
	if err != nil {
		log.Fatalf("Failed to initialize profile store: %v", err)
	}

	h := api.NewHandler(jobStore, profileStore)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/jobs", h.ListJobs)
	mux.HandleFunc("/jobs/scrape", h.ScrapeJobs)
	mux.HandleFunc("/profiles", h.ListProfiles)

	log.Printf("Starting API server on :%s", *port)
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		log.Fatal(err)
	}
}
