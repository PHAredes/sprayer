package scrapers

import (
	"database/sql"
	"log"

	"job-scraper/pkg/models"
)

// Scraper interface defines the contract for all scrapers
type Scraper interface {
	GetName() string
	Scrape() ([]models.Job, error)
}

// ScraperRegistry manages all available scrapers
type ScraperRegistry struct {
	scrapers map[string]Scraper
}

// NewScraperRegistry creates a new registry with all available scrapers
func NewScraperRegistry(db *sql.DB) *ScraperRegistry {
	registry := &ScraperRegistry{
		scrapers: make(map[string]Scraper),
	}

	// Register all scrapers
	registry.scrapers["mock"] = NewMockScraper("mock_scraper", db)
	registry.scrapers["github"] = NewGitHubScraper(db)
	registry.scrapers["ycombinator"] = NewYCombinatorScraper(db)
	registry.scrapers["stackoverflow"] = NewStackOverflowScraper(db)
	registry.scrapers["linkedin"] = NewLinkedInScraper(db)
	registry.scrapers["remote"] = NewRemoteScraper(db)
	registry.scrapers["startup"] = NewStartupScraper(db)

	return registry
}

// GetScraper returns a scraper by name
func (sr *ScraperRegistry) GetScraper(name string) (Scraper, bool) {
	scraper, exists := sr.scrapers[name]
	return scraper, exists
}

// GetAllScrapers returns all registered scrapers
func (sr *ScraperRegistry) GetAllScrapers() []Scraper {
	scrapers := make([]Scraper, 0, len(sr.scrapers))
	for _, scraper := range sr.scrapers {
		scrapers = append(scrapers, scraper)
	}
	return scrapers
}

// GetScraperNames returns the names of all registered scrapers
func (sr *ScraperRegistry) GetScraperNames() []string {
	names := make([]string, 0, len(sr.scrapers))
	for name := range sr.scrapers {
		names = append(names, name)
	}
	return names
}

// ScrapeAll runs all scrapers and returns combined results
func (sr *ScraperRegistry) ScrapeAll() ([]models.Job, error) {
	var allJobs []models.Job
	var totalJobs int

	log.Printf("Starting scraping from %d sources...", len(sr.scrapers))

	for name, scraper := range sr.scrapers {
		log.Printf("Scraping from %s...", name)
		jobs, err := scraper.Scrape()
		if err != nil {
			log.Printf("Warning: Failed to scrape from %s: %v", name, err)
			continue
		}
		allJobs = append(allJobs, jobs...)
		totalJobs += len(jobs)
		log.Printf("Found %d jobs from %s", len(jobs), name)
	}

	log.Printf("Scraping complete! Total jobs found: %d", totalJobs)

	// Check if we reached the target
	if totalJobs >= 300 {
		log.Printf("🎉 Target of 300+ jobs reached!")
	} else {
		log.Printf("Target not reached yet. Current count: %d", totalJobs)
	}

	return allJobs, nil
}