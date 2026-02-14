package scraper

import (
	"fmt"
	"strings"
	"time"

	"sprayer/internal/job"
	"sprayer/internal/parse"

	"github.com/go-rod/rod"
)

// Indeed returns a browser-based scraper for Indeed job search.
func Indeed(query, location string) job.Scraper {
	url := fmt.Sprintf("https://www.indeed.com/jobs?q=%s&l=%s&fromage=7",
		strings.ReplaceAll(query, " ", "+"),
		strings.ReplaceAll(location, " ", "+"),
	)

	return BrowserScrape(url, func(page *rod.Page) ([]job.Job, error) {
		page.MustWaitStable()

		elements, err := page.Elements(".job_seen_beacon, .jobsearch-ResultsList .result")
		if err != nil {
			return nil, fmt.Errorf("Indeed: find cards: %w", err)
		}

		var jobs []job.Job
		for _, el := range elements {
			title, _ := el.Element("h2.jobTitle span, .jobTitle")
			company, _ := el.Element("[data-testid='company-name'], .companyName")
			loc, _ := el.Element("[data-testid='text-location'], .companyLocation")
			link, _ := el.Element("a")
			snippet, _ := el.Element(".job-snippet, .underShelfFooter")

			if title == nil {
				continue
			}

			titleText, _ := title.Text()
			companyText := ""
			if company != nil {
				companyText, _ = company.Text()
			}
			locText := ""
			if loc != nil {
				locText, _ = loc.Text()
			}
			desc := ""
			if snippet != nil {
				desc, _ = snippet.Text()
			}
			href := ""
			if link != nil {
				h, _ := link.Attribute("href")
				if h != nil && *h != "" {
					href = *h
					if !strings.HasPrefix(href, "http") {
						href = "https://www.indeed.com" + href
					}
				}
			}

			j := job.Job{
				ID:          idFromContent("indeed", titleText+companyText),
				Title:       strings.TrimSpace(titleText),
				Company:     strings.TrimSpace(companyText),
				Location:    strings.TrimSpace(locText),
				Description: strings.TrimSpace(desc),
				URL:         href,
				Source:      "indeed",
				PostedDate:  time.Now(),
				Email:       parse.ExtractFirstEmail(desc),
				Salary:      parse.ExtractSalary(desc),
				Score:       50,
			}
			jobs = append(jobs, j)
		}
		return jobs, nil
	})
}

// Glassdoor returns a browser-based scraper for Glassdoor job search.
func Glassdoor(query string) job.Scraper {
	url := fmt.Sprintf("https://www.glassdoor.com/Job/jobs.htm?sc.keyword=%s",
		strings.ReplaceAll(query, " ", "+"),
	)

	return BrowserScrape(url, func(page *rod.Page) ([]job.Job, error) {
		page.MustWaitStable()

		elements, err := page.Elements("[data-test='jobListing'], .react-job-listing")
		if err != nil {
			return nil, fmt.Errorf("Glassdoor: find cards: %w", err)
		}

		var jobs []job.Job
		for _, el := range elements {
			title, _ := el.Element("[data-test='job-title'], .job-title")
			company, _ := el.Element("[data-test='emp-name'], .employer-name")
			loc, _ := el.Element("[data-test='emp-location'], .location")

			if title == nil {
				continue
			}

			titleText, _ := title.Text()
			companyText := ""
			if company != nil {
				companyText, _ = company.Text()
			}
			locText := ""
			if loc != nil {
				locText, _ = loc.Text()
			}

			href := ""
			if link, _ := title.Element("a"); link != nil {
				h, _ := link.Attribute("href")
				if h != nil && *h != "" {
					href = *h
					if !strings.HasPrefix(href, "http") {
						href = "https://www.glassdoor.com" + href
					}
				}
			}

			j := job.Job{
				ID:       idFromContent("gd", titleText+companyText),
				Title:    strings.TrimSpace(titleText),
				Company:  strings.TrimSpace(companyText),
				Location: strings.TrimSpace(locText),
				URL:      href,
				Source:   "glassdoor",
				PostedDate: time.Now(),
				Score:    50,
			}
			jobs = append(jobs, j)
		}
		return jobs, nil
	})
}
