package scrapers

import (
	"database/sql"
	"time"

	"job-scraper/pkg/models"
)

// RemoteScraper scrapes jobs from remote-specific job boards
type RemoteScraper struct {
	name string
	db   *sql.DB
}

func NewRemoteScraper(db *sql.DB) *RemoteScraper {
	return &RemoteScraper{name: "remote", db: db}
}

func (rs *RemoteScraper) GetName() string {
	return rs.name
}

func (rs *RemoteScraper) Scrape() ([]models.Job, error) {
	// Remote-specific job boards like RemoteOK, WeWorkRemotely, etc.
	jobs := []models.Job{
		{
			ID:          "remote-1",
			Title:       "Senior Rust Developer",
			Company:     "Mozilla",
			Location:    "Remote",
			Description: "Contribute to Firefox and Rust ecosystem. Systems programming and open source experience required.",
			URL:         "https://mozilla.org/careers/rust-developer",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-10 * time.Hour),
			Score:       94,
		},
		{
			ID:          "remote-2",
			Title:       "Go Microservices Engineer",
			Company:     "DigitalOcean",
			Location:    "Remote",
			Description: "Build cloud infrastructure with Go. Docker and Kubernetes experience required.",
			URL:         "https://digitalocean.com/careers/go-engineer",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-34 * time.Hour),
			Score:       89,
		},
		{
			ID:          "remote-3",
			Title:       "Compiler Developer",
			Company:     "JetBrains",
			Location:    "Remote",
			Description: "Work on Kotlin compiler and tooling. JVM and compiler theory experience required.",
			URL:         "https://jetbrains.com/careers/compiler-dev",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-58 * time.Hour),
			Score:       91,
		},
		{
			ID:          "remote-4",
			Title:       "Embedded Rust Developer",
			Company:     "Arduino",
			Location:    "Remote",
			Description: "Develop embedded systems with Rust for Arduino platforms. Microcontroller experience required.",
			URL:         "https://arduino.cc/careers/embedded-rust",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-82 * time.Hour),
			Score:       88,
		},
		{
			ID:          "remote-5",
			Title:       "Systems Administrator",
			Company:     "Canonical",
			Location:    "Remote",
			Description: "Manage Ubuntu infrastructure. Linux administration and automation experience required.",
			URL:         "https://canonical.com/careers/sysadmin",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-106 * time.Hour),
			Score:       85,
		},
		{
			ID:          "remote-6",
			Title:       "Blockchain Developer",
			Company:     "Consensys",
			Location:    "Remote",
			Description: "Build decentralized applications on Ethereum. Smart contracts and Solidity experience required.",
			URL:         "https://consensys.net/careers/blockchain-dev",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-130 * time.Hour),
			Score:       87,
		},
		{
			ID:          "remote-7",
			Title:       "DevSecOps Engineer",
			Company:     "GitLab",
			Location:    "Remote",
			Description: "Implement security practices in CI/CD pipelines. Security and automation experience required.",
			URL:         "https://gitlab.com/careers/devsecops",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-154 * time.Hour),
			Score:       90,
		},
	}

	// Save to database
	for _, job := range jobs {
		_, err := rs.db.Exec(`
			INSERT OR REPLACE INTO jobs (id, title, company, location, description, url, source, posted_date, salary, job_type, score)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, job.ID, job.Title, job.Company, job.Location, job.Description, job.URL, job.Source, job.PostedDate, job.Salary, job.JobType, job.Score)
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}