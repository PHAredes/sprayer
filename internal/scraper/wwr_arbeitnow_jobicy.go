package scraper

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"sprayer/internal/job"
	"sprayer/internal/parse"
)

// WeWorkRemotely scrapes the WWR JSON feed.
func WeWorkRemotely() job.Scraper {
	return func() ([]job.Job, error) {
		// WWR exposes category-based JSON feeds
		categories := []string{
			"programming", "devops-sysadmin", "design",
		}

		var all []job.Job
		for _, cat := range categories {
			url := fmt.Sprintf("https://weworkremotely.com/categories/%s/jobs.json", cat)
			data, err := httpGet(url)
			if err != nil {
				continue
			}

			var listings []struct {
				ID          int    `json:"id"`
				Title       string `json:"title"`
				CompanyName string `json:"company_name"`
				Description string `json:"description"`
				URL         string `json:"url"`
				CreatedAt   string `json:"created_at"`
				Category    string `json:"category"`
			}
			if err := json.Unmarshal(data, &listings); err != nil {
				continue
			}

			for _, l := range listings {
				desc := stripHTML(l.Description)
				posted, _ := time.Parse(time.RFC3339, l.CreatedAt)
				j := job.Job{
					ID:          fmt.Sprintf("wwr-%d", l.ID),
					Title:       l.Title,
					Company:     l.CompanyName,
					Location:    "Remote",
					Description: desc,
					URL:         fmt.Sprintf("https://weworkremotely.com%s", l.URL),
					Source:      "weworkremotely",
					PostedDate:  posted,
					Email:       parse.ExtractFirstEmail(desc),
					Salary:      parse.ExtractSalary(desc),
					JobType:     l.Category,
					Score:       50,
				}
				all = append(all, j)
			}
			time.Sleep(200 * time.Millisecond)
		}
		return all, nil
	}
}

// Arbeitnow scrapes the Arbeitnow public JSON API (EU-focused remote jobs).
func Arbeitnow() job.Scraper {
	return func() ([]job.Job, error) {
		var all []job.Job
		page := 1
		for page <= 3 {
			url := fmt.Sprintf("https://www.arbeitnow.com/api/job-board-api?page=%d", page)
			data, err := httpGet(url)
			if err != nil {
				break
			}

			var result struct {
				Data []struct {
					Slug        string   `json:"slug"`
					Title       string   `json:"title"`
					CompanyName string   `json:"company_name"`
					Description string   `json:"description"`
					Location    string   `json:"location"`
					Remote      bool     `json:"remote"`
					URL         string   `json:"url"`
					Tags        []string `json:"tags"`
					CreatedAt   int64    `json:"created_at"`
				} `json:"data"`
				Links struct {
					Next string `json:"next"`
				} `json:"links"`
			}
			if err := json.Unmarshal(data, &result); err != nil {
				break
			}

			for _, d := range result.Data {
				loc := d.Location
				if d.Remote {
					loc = "Remote / " + loc
				}
				desc := stripHTML(d.Description)
				j := job.Job{
					ID:          fmt.Sprintf("an-%s", d.Slug),
					Title:       d.Title,
					Company:     d.CompanyName,
					Location:    loc,
					Description: desc,
					URL:         d.URL,
					Source:      "arbeitnow",
					PostedDate:  time.Unix(d.CreatedAt, 0),
					Email:       parse.ExtractFirstEmail(desc),
					Salary:      parse.ExtractSalary(desc),
					JobType:     strings.Join(d.Tags, ", "),
					Score:       50,
				}
				all = append(all, j)
			}

			if result.Links.Next == "" {
				break
			}
			page++
			time.Sleep(300 * time.Millisecond)
		}
		return all, nil
	}
}

// Jobicy scrapes the Jobicy public API (remote tech jobs).
func Jobicy() job.Scraper {
	return func() ([]job.Job, error) {
		data, err := httpGet("https://jobicy.com/api/v2/remote-jobs?count=50&industry=tech")
		if err != nil {
			return nil, fmt.Errorf("Jobicy API: %w", err)
		}

		var result struct {
			Jobs []struct {
				ID             int    `json:"id"`
				URL            string `json:"url"`
				JobTitle       string `json:"jobTitle"`
				CompanyName    string `json:"companyName"`
				JobGeo         string `json:"jobGeo"`
				JobType        string `json:"jobType"`
				AnnSalaryMin   string `json:"annualSalaryMin"`
				AnnSalaryMax   string `json:"annualSalaryMax"`
				SalaryCurrency string `json:"salaryCurrency"`
				PubDate        string `json:"pubDate"`
				JobExcerpt     string `json:"jobExcerpt"`
			} `json:"jobs"`
		}
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("Jobicy parse: %w", err)
		}

		var jobs []job.Job
		for _, jj := range result.Jobs {
			salary := ""
			if jj.AnnSalaryMin != "" && jj.AnnSalaryMax != "" {
				salary = fmt.Sprintf("%s %s - %s", jj.SalaryCurrency, jj.AnnSalaryMin, jj.AnnSalaryMax)
			}
			posted, _ := time.Parse("2006-01-02 15:04:05", jj.PubDate)
			j := job.Job{
				ID:          fmt.Sprintf("jcy-%d", jj.ID),
				Title:       jj.JobTitle,
				Company:     jj.CompanyName,
				Location:    jj.JobGeo,
				Description: jj.JobExcerpt,
				URL:         jj.URL,
				Source:      "jobicy",
				PostedDate:  posted,
				Salary:      salary,
				JobType:     jj.JobType,
				Score:       50,
			}
			jobs = append(jobs, j)
		}
		return jobs, nil
	}
}
