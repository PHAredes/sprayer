package api

import (
	"encoding/json"
	"log"
	"net/http"

	"sprayer/internal/job"
	"sprayer/internal/profile"
	"sprayer/internal/scraper"
)

type Handler struct {
	store        *job.Store
	profileStore *profile.Store
}

func NewHandler(s *job.Store, p *profile.Store) *Handler {
	return &Handler{store: s, profileStore: p}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": "v1"})
}

func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.store.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Optional filtering query params could be added here
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (h *Handler) ScrapeJobs(w http.ResponseWriter, r *http.Request) {
	// Simple scrape trigger based on query params
	keywords := r.URL.Query()["keywords"]
	if len(keywords) == 0 {
		keywords = []string{"golang", "remote"}
	}

	// API only mode for speed via query param?
	fast := r.URL.Query().Get("fast") == "true"

	var s job.Scraper
	if fast {
		s = scraper.APIOnly()
	} else {
		s = scraper.All(keywords, "Remote")
	}

	go func() {
		jobs, err := s()
		if err != nil {
			log.Printf("scrape error: %v", err)
			return
		}
		if err := h.store.Save(jobs); err != nil {
			log.Printf("save jobs error: %v", err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "scraping started"})
}

func (h *Handler) ListProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.profileStore.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profiles)
}
