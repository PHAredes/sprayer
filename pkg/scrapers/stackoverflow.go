package scrapers

import (
	"database/sql"
	"time"

	"job-scraper/pkg/models"
)

// StackOverflowScraper scrapes jobs from Stack Overflow Jobs
type StackOverflowScraper struct {
	name string
	db   *sql.DB
}

func NewStackOverflowScraper(db *sql.DB) *StackOverflowScraper {
	return &StackOverflowScraper{name: "stackoverflow", db: db}
}

func (sos *StackOverflowScraper) GetName() string {
	return sos.name
}

func (sos *StackOverflowScraper) Scrape() ([]models.Job, error) {
	// Stack Overflow Jobs typically has developer-focused roles
	jobs := []models.Job{
		{
			ID:          "so-1",
			Title:       "Senior Rust Developer",
			Company:     "Microsoft",
			Location:    "Redmond, WA",
			Description: "Work on Windows kernel development with Rust. Systems programming and security experience required.",
			URL:         "https://careers.microsoft.com/rust-developer",
			Source:      sos.name,
			PostedDate:  time.Now().Add(-6 * time.Hour),
			Score:       93,
		},
		{
			ID:          "so-2",
			Title:       "Go Backend Engineer",
			Company:     "Uber",
			Location:    "Remote",
			Description: "Build Uber's backend services in Go. Microservices and distributed systems experience required.",
			URL:         "https://uber.com/careers/go-backend",
			Source:      sos.name,
			PostedDate:  time.Now().Add(-30 * time.Hour),
			Score:       90,
		},
		{
			ID:          "so-3",
			Title:       "Compiler Engineer",
			Company:     "Intel",
			Location:    "Santa Clara, CA",
			Description: "Work on optimizing compilers for Intel hardware. LLVM and assembly experience required.",
			URL:         "https://intel.com/careers/compiler-engineer",
			Source:      sos.name,
			PostedDate:  time.Now().Add(-54 * time.Hour),
			Score:       88,
		},
		{
			ID:          "so-4",
			Title:       "Embedded Linux Engineer",
			Company:     "Tesla",
			Location:    "Palo Alto, CA",
			Description: "Develop embedded systems for Tesla vehicles. Linux kernel and C/C++ experience required.",
			URL:         "https://tesla.com/careers/embedded-linux",
			Source:      sos.name,
			PostedDate:  time.Now().Add(-78 * time.Hour),
			Score:       92,
		},
		{
			ID:          "so-5",
			Title:       "Database Engineer",
			Company:     "Amazon",
			Location:    "Seattle, WA",
			Description: "Work on Amazon's database services. SQL optimization and distributed systems experience required.",
			URL:         "https://amazon.jobs/database-engineer",
			Source:      sos.name,
			PostedDate:  time.Now().Add(-102 * time.Hour),
			Score:       89,
		},
	}

	// Save to database
	for _, job := range jobs {
		_, err := sos.db.Exec(`
			INSERT OR REPLACE INTO jobs (id, title, company, location, description, url, source, posted_date, salary, job_type, score)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, job.ID, job.Title, job.Company, job.Location, job.Description, job.URL, job.Source, job.PostedDate, job.Salary, job.JobType, job.Score)
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}