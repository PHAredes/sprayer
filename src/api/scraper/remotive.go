package scraper

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"sprayer/src/api/job"
	"sprayer/src/api/parse"
)

// Remotive scrapes the Remotive public JSON API.
func Remotive() job.Scraper {
	return func() ([]job.Job, error) {
		// Only fetch software dev jobs to keep it relevant and mostly within limit
		data, err := httpGet("https://remotive.com/api/remote-jobs?category=software-dev")
		if err != nil {
			return nil, fmt.Errorf("Remotive API: %w", err)
		}

		var result struct {
			JobCount int           `json:"job-count"`
			Jobs     []remotiveJob `json:"jobs"`
		}
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("Remotive JSON parse: %w", err)
		}

		var jobs []job.Job
		for _, r := range result.Jobs {
			posted, _ := time.Parse(time.RFC3339, r.PublicationDate)
			if posted.IsZero() {
				posted = time.Now()
			}

			desc := stripHTML(r.Description)

			j := job.Job{
				ID:          fmt.Sprintf("remotive-%d", r.ID),
				Title:       r.Title,
				Company:     r.CompanyName,
				Location:    r.CandidateRequiredLocation,
				Description: desc,
				URL:         r.URL,
				Source:      "remotive",
				PostedDate:  posted,
				Email:       parse.ExtractFirstEmail(desc),
				Salary:      r.Salary,
				JobType:     strings.Join(r.Tags, ", "),
				Score:       50,
			}
			jobs = append(jobs, j)
		}

		return jobs, nil
	}
}

type remotiveJob struct {
	ID                        int      `json:"id"`
	URL                       string   `json:"url"`
	Title                     string   `json:"title"`
	CompanyName               string   `json:"company_name"`
	Category                  string   `json:"category"`
	Tags                      []string `json:"tags"`
	JobType                   string   `json:"job_type"`
	PublicationDate           string   `json:"publication_date"`
	CandidateRequiredLocation string   `json:"candidate_required_location"`
	Salary                    string   `json:"salary"`
	Description               string   `json:"description"`
}
