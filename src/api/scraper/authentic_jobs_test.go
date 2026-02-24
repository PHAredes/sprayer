package scraper

import (
	"testing"
	"time"
)

func TestAuthenticJobs(t *testing.T) {
	scraper := AuthenticJobs()

	// Test that the scraper can fetch and parse jobs
	jobs, err := scraper()
	if err != nil {
		t.Fatalf("AuthenticJobs scraper failed: %v", err)
	}

	if len(jobs) == 0 {
		t.Fatal("Expected at least one job from AuthenticJobs feed")
	}

	// Test the first job has required fields
	job := jobs[0]
	if job.Title == "" {
		t.Error("Job title should not be empty")
	}
	if job.Company == "" {
		t.Error("Job company should not be empty")
	}
	if job.URL == "" {
		t.Error("Job URL should not be empty")
	}
	if job.Source != "authenticjobs" {
		t.Errorf("Expected source 'authenticjobs', got '%s'", job.Source)
	}
	if job.PostedDate.IsZero() {
		t.Error("Job posted date should not be zero")
	}
	if job.Description == "" {
		t.Error("Job description should not be empty")
	}

	// Log some details for verification
	t.Logf("Found %d jobs from AuthenticJobs", len(jobs))
	t.Logf("First job: %s at %s (%s)", job.Title, job.Company, job.Location)
	t.Logf("Job type: %s", job.JobType)
	t.Logf("Posted: %s", job.PostedDate.Format(time.RFC3339))
}
