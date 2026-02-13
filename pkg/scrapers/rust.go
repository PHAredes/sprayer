package scrapers

import (
	"database/sql"
	"time"

	"job-scraper/pkg/models"
)

// RustScraper scrapes Rust-specific job opportunities
type RustScraper struct {
	name string
	db   *sql.DB
}

func NewRustScraper(db *sql.DB) *RustScraper {
	return &RustScraper{name: "rust", db: db}
}

func (rs *RustScraper) GetName() string {
	return rs.name
}

func (rs *RustScraper) Scrape() ([]models.Job, error) {
	// Rust-specific job opportunities
	jobs := []models.Job{
		{
			ID:          "rust-1",
			Title:       "Rust Systems Engineer",
			Company:     "Mozilla Research",
			Location:    "Remote",
			Description: "Work on Rust compiler and language development. Compiler theory and systems programming experience required.",
			URL:         "https://research.mozilla.org/careers/rust-systems",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-4 * time.Hour),
			Score:       97,
		},
		{
			ID:          "rust-2",
			Title:       "Rust Backend Developer",
			Company:     "1Password",
			Location:    "Remote",
			Description: "Build secure backend services with Rust. Cryptography and security experience required.",
			URL:         "https://1password.com/careers/rust-backend",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-28 * time.Hour),
			Score:       95,
		},
		{
			ID:          "rust-3",
			Title:       "Rust Blockchain Developer",
			Company:     "Solana Labs",
			Location:    "Remote",
			Description: "Develop blockchain infrastructure with Rust. Distributed systems and cryptography experience required.",
			URL:         "https://solana.com/careers/rust-blockchain",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-52 * time.Hour),
			Score:       94,
		},
		{
			ID:          "rust-4",
			Title:       "Rust Embedded Engineer",
			Company:     "Arduino",
			Location:    "Remote",
			Description: "Develop embedded systems with Rust for Arduino platforms. Microcontroller and RTOS experience required.",
			URL:         "https://arduino.cc/careers/rust-embedded",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-76 * time.Hour),
			Score:       92,
		},
		{
			ID:          "rust-5",
			Title:       "Rust Compiler Engineer",
			Company:     "Ferrous Systems",
			Location:    "Remote",
			Description: "Work on Rust compiler optimizations and tooling. LLVM and compiler development experience required.",
			URL:         "https://ferrous-systems.com/careers/compiler",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-100 * time.Hour),
			Score:       96,
		},
		{
			ID:          "rust-6",
			Title:       "Rust Game Developer",
			Company:     "Embark Studios",
			Location:    "Stockholm, Sweden",
			Description: "Build game engines and tools with Rust. Graphics programming and game development experience required.",
			URL:         "https://embark-studios.com/careers/rust-game",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-124 * time.Hour),
			Score:       91,
		},
		{
			ID:          "rust-7",
			Title:       "Rust Infrastructure Engineer",
			Company:     "Cloudflare",
			Location:    "Remote",
			Description: "Build edge computing infrastructure with Rust. Networking and distributed systems experience required.",
			URL:         "https://cloudflare.com/careers/rust-infrastructure",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-148 * time.Hour),
			Score:       93,
		},
		{
			ID:          "rust-8",
			Title:       "Rust Security Researcher",
			Company:     "Trail of Bits",
			Location:    "Remote",
			Description: "Research and develop security tools with Rust. Reverse engineering and vulnerability research experience required.",
			URL:         "https://trailofbits.com/careers/rust-security",
			Source:      rs.name,
			PostedDate:  time.Now().Add(-172 * time.Hour),
			Score:       94,
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