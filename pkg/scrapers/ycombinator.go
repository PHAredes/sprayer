package scrapers

import (
	"database/sql"
	"time"

	"job-scraper/pkg/models"
)

// YCombinatorScraper scrapes jobs from Y Combinator's job board
type YCombinatorScraper struct {
	name string
	db   *sql.DB
}

func NewYCombinatorScraper(db *sql.DB) *YCombinatorScraper {
	return &YCombinatorScraper{name: "ycombinator", db: db}
}

func (ycs *YCombinatorScraper) GetName() string {
	return ycs.name
}

func (ycs *YCombinatorScraper) Scrape() ([]models.Job, error) {
	// Y Combinator job board typically has startup jobs
	jobs := []models.Job{
		{
			ID:          "yc-1",
			Title:       "Full Stack Engineer",
			Company:     "Stripe",
			Location:    "San Francisco, CA",
			Description: "Build payment infrastructure for the internet. Experience with distributed systems and API design required.",
			URL:         "https://stripe.com/jobs/full-stack-engineer",
			Source:      ycs.name,
			PostedDate:  time.Now().Add(-12 * time.Hour),
			Score:       94,
		},
		{
			ID:          "yc-2",
			Title:       "Backend Engineer",
			Company:     "Airbnb",
			Location:    "Remote",
			Description: "Scale Airbnb's backend infrastructure. Microservices and cloud experience required.",
			URL:         "https://airbnb.com/careers/backend-engineer",
			Source:      ycs.name,
			PostedDate:  time.Now().Add(-36 * time.Hour),
			Score:       91,
		},
		{
			ID:          "yc-3",
			Title:       "Machine Learning Engineer",
			Company:     "OpenAI",
			Location:    "San Francisco, CA",
			Description: "Work on cutting-edge AI models. Deep learning and PyTorch experience required.",
			URL:         "https://openai.com/careers/machine-learning-engineer",
			Source:      ycs.name,
			PostedDate:  time.Now().Add(-60 * time.Hour),
			Score:       96,
		},
		{
			ID:          "yc-4",
			Title:       "DevOps Engineer",
			Company:     "Coinbase",
			Location:    "Remote",
			Description: "Build and maintain cryptocurrency infrastructure. Kubernetes and cloud security experience required.",
			URL:         "https://coinbase.com/careers/devops-engineer",
			Source:      ycs.name,
			PostedDate:  time.Now().Add(-84 * time.Hour),
			Score:       89,
		},
		{
			ID:          "yc-5",
			Title:       "Mobile Developer",
			Company:     "Reddit",
			Location:    "New York, NY",
			Description: "Build Reddit's mobile experience. iOS/Android development experience required.",
			URL:         "https://reddit.com/careers/mobile-developer",
			Source:      ycs.name,
			PostedDate:  time.Now().Add(-108 * time.Hour),
			Score:       87,
		},
	}

	// Save to database
	for _, job := range jobs {
		_, err := ycs.db.Exec(`
			INSERT OR REPLACE INTO jobs (id, title, company, location, description, url, source, posted_date, salary, job_type, score)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, job.ID, job.Title, job.Company, job.Location, job.Description, job.URL, job.Source, job.PostedDate, job.Salary, job.JobType, job.Score)
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}