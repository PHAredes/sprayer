package scraper

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"sprayer/internal/job"
	"sprayer/internal/parse"
)

// DefaultGreenhouseBoards is a curated list of companies using Greenhouse.
var DefaultGreenhouseBoards = []string{
	"flyio", "cloudflare", "cockroachlabs", "figma",
	"notion", "vercel", "planetscale", "linear",
}

// Greenhouse scrapes the Greenhouse JSON API for a set of company boards.
func Greenhouse(boards []string) job.Scraper {
	return func() ([]job.Job, error) {
		var all []job.Job
		for _, board := range boards {
			jobs, err := scrapeGreenhouseBoard(board)
			if err != nil {
				continue // Skip failing boards
			}
			all = append(all, jobs...)
			time.Sleep(300 * time.Millisecond) // Rate limit
		}
		return all, nil
	}
}

func scrapeGreenhouseBoard(board string) ([]job.Job, error) {
	url := fmt.Sprintf("https://boards-api.greenhouse.io/v1/boards/%s/jobs?content=true", board)
	data, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	var result struct {
		Jobs []greenhouseJob `json:"jobs"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	var jobs []job.Job
	for _, gj := range result.Jobs {
		loc := gj.Location.Name
		desc := stripHTML(gj.Content)

		j := job.Job{
			ID:          fmt.Sprintf("gh-%s-%d", board, gj.ID),
			Title:       gj.Title,
			Company:     board,
			Location:    loc,
			Description: desc,
			URL:         gj.AbsoluteURL,
			Source:      "greenhouse",
			PostedDate:  gj.UpdatedAt,
			Email:       parse.ExtractFirstEmail(desc),
			Salary:      parse.ExtractSalary(desc),
			Score:       50,
		}
		jobs = append(jobs, j)
	}

	return jobs, nil
}

type greenhouseJob struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	AbsoluteURL string   `json:"absolute_url"`
	UpdatedAt   time.Time `json:"updated_at"`
	Location    struct {
		Name string `json:"name"`
	} `json:"location"`
	Departments []struct {
		Name string `json:"name"`
	} `json:"departments"`
}

// GreenhouseForKeywords returns jobs filtered by keywords in title.
func GreenhouseForKeywords(boards []string, keywords []string) job.Scraper {
	base := Greenhouse(boards)
	filter := job.ByKeywords(keywords)
	return func() ([]job.Job, error) {
		jobs, err := base()
		if err != nil {
			return nil, err
		}
		return filter(jobs), nil
	}
}

// CompanyNameFromBoard prettifies a board slug.
func CompanyNameFromBoard(board string) string {
	replacer := strings.NewReplacer(
		"flyio", "Fly.io",
		"cloudflare", "Cloudflare",
		"cockroachlabs", "Cockroach Labs",
		"figma", "Figma",
		"notion", "Notion",
		"vercel", "Vercel",
		"planetscale", "PlanetScale",
		"linear", "Linear",
	)
	name := replacer.Replace(board)
	if name == board {
		return strings.Title(board)
	}
	return name
}
