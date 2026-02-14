package scraper

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"sprayer/internal/job"
)

// RemoteOK scrapes the RemoteOK public JSON API.
func RemoteOK() job.Scraper {
	return func() ([]job.Job, error) {
		data, err := httpGet("https://remoteok.com/api")
		if err != nil {
			return nil, fmt.Errorf("RemoteOK API: %w", err)
		}

		// RemoteOK returns an array where [0] is metadata, rest are jobs
		var raw []json.RawMessage
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, fmt.Errorf("RemoteOK parse: %w", err)
		}

		var jobs []job.Job
		for i, entry := range raw {
			if i == 0 {
				continue // Skip metadata entry
			}

			var r remoteOKJob
			if err := json.Unmarshal(entry, &r); err != nil {
				continue
			}

			posted := time.Unix(r.Epoch, 0)

			j := job.Job{
				ID:          fmt.Sprintf("rok-%s", r.ID),
				Title:       r.Position,
				Company:     r.Company,
				Location:    strings.Join(r.Location, ", "),
				Description: stripHTML(r.Description),
				URL:         fmt.Sprintf("https://remoteok.com/remote-jobs/%s", r.Slug),
				Source:      "remoteok",
				PostedDate:  posted,
				Salary:      formatSalary(r.SalaryMin, r.SalaryMax),
				JobType:     strings.Join(r.Tags, ", "),
				Score:       50,
			}
			jobs = append(jobs, j)
		}

		return jobs, nil
	}
}

type remoteOKJob struct {
	ID          string   `json:"id"`
	Slug        string   `json:"slug"`
	Position    string   `json:"position"`
	Company     string   `json:"company"`
	Description string   `json:"description"`
	Location    []string `json:"location"`
	Tags        []string `json:"tags"`
	SalaryMin   int      `json:"salary_min"`
	SalaryMax   int      `json:"salary_max"`
	Epoch       int64    `json:"epoch"`
}

func formatSalary(min, max int) string {
	if min <= 0 && max <= 0 {
		return ""
	}
	if min > 0 && max > 0 {
		return fmt.Sprintf("$%dk - $%dk", min/1000, max/1000)
	}
	if min > 0 {
		return fmt.Sprintf("$%dk+", min/1000)
	}
	return fmt.Sprintf("up to $%dk", max/1000)
}
