package tracking

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// TrackingServer handles HTTP requests for tracking pixels and links
type TrackingServer struct {
	Port int
}

// NewTrackingServer creates a new tracking server
func NewTrackingServer(port int) *TrackingServer {
	return &TrackingServer{
		Port: port,
	}
}

// Start starts the tracking server
func (ts *TrackingServer) Start() error {
	// Set up routes
	http.HandleFunc("/track/open/", ts.handleEmailOpen)
	http.HandleFunc("/track/click/", ts.handleLinkClick)
	http.HandleFunc("/health", ts.handleHealth)

	addr := fmt.Sprintf(":%d", ts.Port)
	log.Printf("Starting tracking server on port %d", ts.Port)
	return http.ListenAndServe(addr, nil)
}

// handleEmailOpen handles email open tracking requests
func (ts *TrackingServer) handleEmailOpen(w http.ResponseWriter, r *http.Request) {
	// Extract job ID and email hash from URL path
	// Expected format: /track/open/{jobID}/{emailHash}.gif
	path := strings.TrimPrefix(r.URL.Path, "/track/open/")
	parts := strings.Split(path, "/")
	
	if len(parts) < 2 {
		http.Error(w, "Invalid tracking URL", http.StatusBadRequest)
		return
	}

	jobID := parts[0]
	emailHash := strings.TrimSuffix(parts[1], ".gif")

	TrackEmailOpen(w, r, jobID, emailHash)
}

// handleLinkClick handles link click tracking requests
func (ts *TrackingServer) handleLinkClick(w http.ResponseWriter, r *http.Request) {
	// Extract job ID and email hash from URL path
	// Expected format: /track/click/{jobID}/{emailHash}
	path := strings.TrimPrefix(r.URL.Path, "/track/click/")
	parts := strings.Split(path, "/")
	
	if len(parts) < 2 {
		http.Error(w, "Invalid tracking URL", http.StatusBadRequest)
		return
	}

	jobID := parts[0]
	emailHash := parts[1]

	TrackLinkClick(w, r, jobID, emailHash)
}

// handleHealth handles health check requests
func (ts *TrackingServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok", "service": "tracking"}`))
}