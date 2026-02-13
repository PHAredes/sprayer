package scrapers

import (
	"database/sql"
	"time"

	"job-scraper/pkg/models"
)

// LinkedInScraper scrapes jobs from LinkedIn (simulated)
type LinkedInScraper struct {
	name string
	db   *sql.DB
}

func NewLinkedInScraper(db *sql.DB) *LinkedInScraper {
	return &LinkedInScraper{name: "linkedin", db: db}
}

func (ls *LinkedInScraper) GetName() string {
	return ls.name
}

func (ls *LinkedInScraper) Scrape() ([]models.Job, error) {
	// LinkedIn has diverse job postings across industries
	jobs := []models.Job{
		{
			ID:          "li-1",
			Title:       "Systems Software Engineer",
			Company:     "Apple",
			Location:    "Cupertino, CA",
			Description: "Work on macOS and iOS kernel development. C/C++ and assembly experience required.",
			URL:         "https://apple.com/careers/systems-software",
			Source:      ls.name,
			PostedDate:  time.Now().Add(-8 * time.Hour),
			Score:       95,
		},
		{
			ID:          "li-2",
			Title:       "Cloud Infrastructure Engineer",
			Company:     "Netflix",
			Location:    "Los Gatos, CA",
			Description: "Build and scale Netflix's cloud infrastructure. AWS and containerization experience required.",
			URL:         "https://jobs.netflix.com/cloud-infrastructure",
			Source:      ls.name,
			PostedDate:  time.Now().Add(-32 * time.Hour),
			Score:       91,
		},
		{
			ID:          "li-3",
			Title:       "Security Researcher",
			Company:     "Meta",
			Location:    "Menlo Park, CA",
			Description: "Research and develop security solutions for Meta's platforms. Reverse engineering experience required.",
			URL:         "https://meta.com/careers/security-researcher",
			Source:      ls.name,
			PostedDate:  time.Now().Add(-56 * time.Hour),
			Score:       90,
		},
		{
			ID:          "li-4",
			Title:       "Performance Engineer",
			Company:     "Oracle",
			Location:    "Austin, TX",
			Description: "Optimize database performance. SQL tuning and benchmarking experience required.",
			URL:         "https://oracle.com/careers/performance-engineer",
			Source:      ls.name,
			PostedDate:  time.Now().Add(-80 * time.Hour),
			Score:       86,
		},
		{
			ID:          "li-5",
			Title:       "Full Stack Developer",
			Company:     "Salesforce",
			Location:    "San Francisco, CA",
			Description: "Build Salesforce's customer-facing applications. JavaScript and cloud experience required.",
			URL:         "https://salesforce.com/careers/full-stack",
			Source:      ls.name,
			PostedDate:  time.Now().Add(-104 * time.Hour),
			Score:       87,
		},
	}

	// Save to database
	for _, job := range jobs {
		_, err := ls.db.Exec(`
			INSERT OR REPLACE INTO jobs (id, title, company, location, description, url, source, posted_date, salary, job_type, score)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, job.ID, job.Title, job.Company, job.Location, job.Description, job.URL, job.Source, job.PostedDate, job.Salary, job.JobType, job.Score)
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}