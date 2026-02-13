package scrapers

import (
	"database/sql"
	"time"

	"job-scraper/pkg/models"
)

// GitHubScraper scrapes jobs from GitHub Jobs API
type GitHubScraper struct {
	name string
	db   *sql.DB
}

func NewGitHubScraper(db *sql.DB) *GitHubScraper {
	return &GitHubScraper{name: "github", db: db}
}

func (ghs *GitHubScraper) GetName() string {
	return ghs.name
}

func (ghs *GitHubScraper) Scrape() ([]models.Job, error) {
	// GitHub Jobs API is deprecated, so we'll use GitHub's own jobs page
	// For now, return mock data simulating GitHub job postings
	jobs := []models.Job{
		{
			ID:          "github-1",
			Title:       "Senior Software Engineer - Infrastructure",
			Company:     "GitHub",
			Location:    "Remote",
			Description: "Build and scale GitHub's infrastructure. Experience with distributed systems and Go/Rust required.",
			URL:         "https://github.com/careers/senior-software-engineer-infrastructure",
			Source:      ghs.name,
			PostedDate:  time.Now().Add(-24 * time.Hour),
			Score:       92,
		},
		{
			ID:          "github-2",
			Title:       "Developer Advocate",
			Company:     "GitHub",
			Location:    "San Francisco, CA",
			Description: "Help developers succeed with GitHub. Technical writing and community experience required.",
			URL:         "https://github.com/careers/developer-advocate",
			Source:      ghs.name,
			PostedDate:  time.Now().Add(-48 * time.Hour),
			Score:       85,
		},
		{
			ID:          "github-3",
			Title:       "Security Engineer",
			Company:     "GitHub",
			Location:    "Remote",
			Description: "Protect GitHub's platform and users. Security research and code review experience required.",
			URL:         "https://github.com/careers/security-engineer",
			Source:      ghs.name,
			PostedDate:  time.Now().Add(-72 * time.Hour),
			Score:       90,
		},
		{
			ID:          "github-4",
			Title:       "Frontend Engineer",
			Company:     "GitHub",
			Location:    "Remote",
			Description: "Build GitHub's web interface. React, TypeScript, and modern web development experience required.",
			URL:         "https://github.com/careers/frontend-engineer",
			Source:      ghs.name,
			PostedDate:  time.Now().Add(-96 * time.Hour),
			Score:       88,
		},
	}

	// Save to database
	for _, job := range jobs {
		_, err := ghs.db.Exec(`
			INSERT OR REPLACE INTO jobs (id, title, company, location, description, url, source, posted_date, salary, job_type, score)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, job.ID, job.Title, job.Company, job.Location, job.Description, job.URL, job.Source, job.PostedDate, job.Salary, job.JobType, job.Score)
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}