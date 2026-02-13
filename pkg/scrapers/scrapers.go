package scrapers

import (
	"database/sql"
	"time"

	"job-scraper/pkg/models"
)

// MockScraper returns test data
type MockScraper struct {
	name string
	db   *sql.DB
}

func NewMockScraper(name string, db *sql.DB) *MockScraper {
	return &MockScraper{name: name, db: db}
}

func (ms *MockScraper) GetName() string {
	return ms.name
}

func (ms *MockScraper) Scrape() ([]models.Job, error) {
	// Check if we already have jobs in database
	var count int
	err := ms.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE source = ?", ms.name).Scan(&count)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		// Return jobs from database
		rows, err := ms.db.Query("SELECT id, title, company, location, description, url, source, posted_date, salary, job_type, score FROM jobs WHERE source = ?", ms.name)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var jobs []models.Job
		for rows.Next() {
			var job models.Job
			err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Location, &job.Description, &job.URL, &job.Source, &job.PostedDate, &job.Salary, &job.JobType, &job.Score)
			if err != nil {
				return nil, err
			}
			jobs = append(jobs, job)
		}
		return jobs, nil
	}

	// Return mock job data and save to database
	jobs := []models.Job{
		{
			ID:          "1",
			Title:       "Senior Rust Developer",
			Company:     "SpaceX",
			Location:    "Hawthorne, CA",
			Description: "Looking for experienced Rust developers to work on Starshield systems. Requirements: 5+ years Rust experience, systems programming background.",
			URL:         "https://spacex.com/careers/rust-dev",
			Source:      ms.name,
			PostedDate:  time.Now(),
			Score:       95,
		},
		{
			ID:          "2",
			Title:       "Compiler Engineer",
			Company:     "Google",
			Location:    "Mountain View, CA",
			Description: "Work on compilers, runtimes and toolchains for Google's infrastructure. LLVM experience preferred.",
			URL:         "https://google.com/careers/compiler",
			Source:      ms.name,
			PostedDate:  time.Now().Add(-24 * time.Hour),
			Score:       90,
		},
		{
			ID:          "3",
			Title:       "Embedded Systems Engineer",
			Company:     "NVIDIA",
			Location:    "Santa Clara, CA",
			Description: "Work on GPU and SOC firmware development. C/C++ and assembly experience required.",
			URL:         "https://nvidia.com/careers/embedded",
			Source:      ms.name,
			PostedDate:  time.Now().Add(-48 * time.Hour),
			Score:       88,
		},
		{
			ID:          "4",
			Title:       "Systems Programmer",
			Company:     "Jito Labs",
			Location:    "Remote",
			Description: "Low latency networking and systems programming in Rust. Blockchain experience a plus.",
			URL:         "https://jito.network/careers",
			Source:      ms.name,
			PostedDate:  time.Now().Add(-72 * time.Hour),
			Score:       92,
		},
		{
			ID:          "5",
			Title:       "Open Source Compiler Developer",
			Company:     "IBM",
			Location:    "Remote",
			Description: "Contribute to open source compiler projects. GCC and LLVM experience preferred.",
			URL:         "https://ibm.com/careers/compiler",
			Source:      ms.name,
			PostedDate:  time.Now().Add(-96 * time.Hour),
			Score:       85,
		},
	}

	// Save jobs to database
	for _, job := range jobs {
		_, err := ms.db.Exec(`
			INSERT OR REPLACE INTO jobs (id, title, company, location, description, url, source, posted_date, salary, job_type, score)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, job.ID, job.Title, job.Company, job.Location, job.Description, job.URL, job.Source, job.PostedDate, job.Salary, job.JobType, job.Score)
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}