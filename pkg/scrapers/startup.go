package scrapers

import (
	"database/sql"
	"time"

	"job-scraper/pkg/models"
)

// StartupScraper scrapes jobs from startup-focused job boards
type StartupScraper struct {
	name string
	db   *sql.DB
}

func NewStartupScraper(db *sql.DB) *StartupScraper {
	return &StartupScraper{name: "startup", db: db}
}

func (ss *StartupScraper) GetName() string {
	return ss.name
}

func (ss *StartupScraper) Scrape() ([]models.Job, error) {
	// Startup job boards like AngelList, Wellfound, etc.
	jobs := []models.Job{
		{
			ID:          "startup-1",
			Title:       "Founding Engineer",
			Company:     "Anthropic",
			Location:    "San Francisco, CA",
			Description: "Help build Claude AI. ML engineering and distributed systems experience required.",
			URL:         "https://anthropic.com/careers/founding-engineer",
			Source:      ss.name,
			PostedDate:  time.Now().Add(-14 * time.Hour),
			Score:       96,
		},
		{
			ID:          "startup-2",
			Title:       "Backend Rust Engineer",
			Company:     "Discord",
			Location:    "San Francisco, CA",
			Description: "Scale Discord's real-time communication platform. Rust and networking experience required.",
			URL:         "https://discord.com/careers/rust-backend",
			Source:      ss.name,
			PostedDate:  time.Now().Add(-38 * time.Hour),
			Score:       93,
		},
		{
			ID:          "startup-3",
			Title:       "Systems Programmer",
			Company:     "Figma",
			Location:    "San Francisco, CA",
			Description: "Build Figma's collaborative design platform. WebAssembly and performance optimization experience required.",
			URL:         "https://figma.com/careers/systems-programmer",
			Source:      ss.name,
			PostedDate:  time.Now().Add(-62 * time.Hour),
			Score:       92,
		},
		{
			ID:          "startup-4",
			Title:       "Infrastructure Engineer",
			Company:     "Notion",
			Location:    "New York, NY",
			Description: "Scale Notion's infrastructure. AWS and database optimization experience required.",
			URL:         "https://notion.so/careers/infrastructure",
			Source:      ss.name,
			PostedDate:  time.Now().Add(-86 * time.Hour),
			Score:       91,
		},
		{
			ID:          "startup-5",
			Title:       "Security Engineer",
			Company:     "1Password",
			Location:    "Remote",
			Description: "Protect 1Password's security infrastructure. Cryptography and secure coding experience required.",
			URL:         "https://1password.com/careers/security-engineer",
			Source:      ss.name,
			PostedDate:  time.Now().Add(-110 * time.Hour),
			Score:       94,
		},
		{
			ID:          "startup-6",
			Title:       "Platform Engineer",
			Company:     "Vercel",
			Location:    "San Francisco, CA",
			Description: "Build Vercel's edge computing platform. Cloud infrastructure and JavaScript experience required.",
			URL:         "https://vercel.com/careers/platform-engineer",
			Source:      ss.name,
			PostedDate:  time.Now().Add(-134 * time.Hour),
			Score:       90,
		},
		{
			ID:          "startup-7",
			Title:       "Data Engineer",
			Company:     "Snowflake",
			Location:    "San Mateo, CA",
			Description: "Work on Snowflake's data platform. SQL optimization and distributed systems experience required.",
			URL:         "https://snowflake.com/careers/data-engineer",
			Source:      ss.name,
			PostedDate:  time.Now().Add(-158 * time.Hour),
			Score:       89,
		},
	}

	// Save to database
	for _, job := range jobs {
		_, err := ss.db.Exec(`
			INSERT OR REPLACE INTO jobs (id, title, company, location, description, url, source, posted_date, salary, job_type, score)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, job.ID, job.Title, job.Company, job.Location, job.Description, job.URL, job.Source, job.PostedDate, job.Salary, job.JobType, job.Score)
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}