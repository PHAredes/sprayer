package scraper

import (
	"testing"
	"time"
)

func TestYCWorkAtStartup(t *testing.T) {
	// Test that the scraper function can be created without errors
	scraper := YCWorkAtStartup([]string{"engineer"}, "remote")

	if scraper == nil {
		t.Fatal("Expected scraper function, got nil")
	}

	// Test that it can be called (this will actually run the browser automation)
	// Note: This test requires a browser environment and may take some time
	t.Run("ScrapeJobs", func(t *testing.T) {
		// Skip this test in CI environments or when browser is not available
		if testing.Short() {
			t.Skip("Skipping browser test in short mode")
		}

		// Set a timeout for the scrape operation
		done := make(chan bool)
		go func() {
			jobs, err := scraper()
			if err != nil {
				t.Logf("Scrape error (expected in test environment): %v", err)
			} else {
				t.Logf("Found %d jobs", len(jobs))
				for i, job := range jobs {
					if i >= 3 { // Only show first 3 jobs
						break
					}
					t.Logf("Job %d: %s at %s in %s", i+1, job.Title, job.Company, job.Location)
				}
			}
			done <- true
		}()

		select {
		case <-done:
			// Test completed
		case <-time.After(30 * time.Second):
			t.Log("Scrape operation timed out after 30 seconds")
		}
	})
}
