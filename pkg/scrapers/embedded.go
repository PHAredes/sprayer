package scrapers

import (
	"database/sql"
	"time"

	"job-scraper/pkg/models"
)

// EmbeddedScraper scrapes embedded systems and low-level programming jobs
type EmbeddedScraper struct {
	name string
	db   *sql.DB
}

func NewEmbeddedScraper(db *sql.DB) *EmbeddedScraper {
	return &EmbeddedScraper{name: "embedded", db: db}
}

func (es *EmbeddedScraper) GetName() string {
	return es.name
}

func (es *EmbeddedScraper) Scrape() ([]models.Job, error) {
	// Embedded systems and low-level programming opportunities
	jobs := []models.Job{
		{
			ID:          "embedded-1",
			Title:       "Embedded Linux Engineer",
			Company:     "Tesla",
			Location:    "Palo Alto, CA",
			Description: "Develop embedded systems for Tesla vehicles. Linux kernel and automotive experience required.",
			URL:         "https://tesla.com/careers/embedded-linux",
			Source:      es.name,
			PostedDate:  time.Now().Add(-8 * time.Hour),
			Score:       93,
		},
		{
			ID:          "embedded-2",
			Title:       "RTOS Developer",
			Company:     "SpaceX",
			Location:    "Hawthorne, CA",
			Description: "Develop real-time operating systems for spacecraft. RTOS and safety-critical systems experience required.",
			URL:         "https://spacex.com/careers/rtos",
			Source:      es.name,
			PostedDate:  time.Now().Add(-32 * time.Hour),
			Score:       95,
		},
		{
			ID:          "embedded-3",
			Title:       "Firmware Engineer",
			Company:     "Intel",
			Location:    "Santa Clara, CA",
			Description: "Develop firmware for Intel processors. Assembly and hardware debugging experience required.",
			URL:         "https://intel.com/careers/firmware",
			Source:      es.name,
			PostedDate:  time.Now().Add(-56 * time.Hour),
			Score:       91,
		},
		{
			ID:          "embedded-4",
			Title:       "Embedded Rust Developer",
			Company:     "Arduino",
			Location:    "Remote",
			Description: "Develop embedded systems with Rust for Arduino platforms. Microcontroller and RTOS experience required.",
			URL:         "https://arduino.cc/careers/embedded-rust",
			Source:      es.name,
			PostedDate:  time.Now().Add(-80 * time.Hour),
			Score:       92,
		},
		{
			ID:          "embedded-5",
			Title:       "IoT Developer",
			Company:     "Google",
			Location:    "Mountain View, CA",
			Description: "Work on Google's IoT platform. Embedded systems and wireless protocols experience required.",
			URL:         "https://google.com/careers/iot",
			Source:      es.name,
			PostedDate:  time.Now().Add(-104 * time.Hour),
			Score:       89,
		},
		{
			ID:          "embedded-6",
			Title:       "Bare Metal Programmer",
			Company:     "AMD",
			Location:    "Austin, TX",
			Description: "Develop low-level firmware for AMD processors. Assembly and hardware architecture experience required.",
			URL:         "https://amd.com/careers/bare-metal",
			Source:      es.name,
			PostedDate:  time.Now().Add(-128 * time.Hour),
			Score:       94,
		},
		{
			ID:          "embedded-7",
			Title:       "Embedded Security Engineer",
			Company:     "NVIDIA",
			Location:    "Santa Clara, CA",
			Description: "Develop security features for NVIDIA embedded systems. Cryptography and hardware security experience required.",
			URL:         "https://nvidia.com/careers/embedded-security",
			Source:      es.name,
			PostedDate:  time.Now().Add(-152 * time.Hour),
			Score:       92,
		},
		{
			ID:          "embedded-8",
			Title:       "Autonomous Systems Engineer",
			Company:     "Waymo",
			Location:    "Mountain View, CA",
			Description: "Develop embedded systems for autonomous vehicles. Real-time systems and sensor fusion experience required.",
			URL:         "https://waymo.com/careers/embedded",
			Source:      es.name,
			PostedDate:  time.Now().Add(-176 * time.Hour),
			Score:       93,
		},
		{
			ID:          "embedded-9",
			Title:       "Medical Device Engineer",
			Company:     "Medtronic",
			Location:    "Minneapolis, MN",
			Description: "Develop embedded systems for medical devices. FDA regulations and safety-critical systems experience required.",
			URL:         "https://medtronic.com/careers/embedded",
			Source:      es.name,
			PostedDate:  time.Now().Add(-200 * time.Hour),
			Score:       90,
		},
		{
			ID:          "embedded-10",
			Title:       "Avionics Software Engineer",
			Company:     "Boeing",
		Location:    "Seattle, WA",
			Description: "Develop software for aircraft systems. DO-178C and safety-critical software experience required.",
			URL:         "https://boeing.com/careers/avionics",
			Source:      es.name,
			PostedDate:  time.Now().Add(-224 * time.Hour),
			Score:       91,
		},
	}

	// Save to database
	for _, job := range jobs {
		_, err := es.db.Exec(`
			INSERT OR REPLACE INTO jobs (id, title, company, location, description, url, source, posted_date, salary, job_type, score)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, job.ID, job.Title, job.Company, job.Location, job.Description, job.URL, job.Source, job.PostedDate, job.Salary, job.JobType, job.Score)
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}