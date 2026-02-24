package scraper

import (
	"fmt"
	"strings"
	"time"

	"sprayer/src/api/job"
	"sprayer/src/api/parse"

	"github.com/go-rod/rod"
)

// Dice returns a browser-based scraper for Dice.com job search.
// It navigates through multiple pages (/jobs/pages/1,2,3...) and extracts job details.
func Dice(keywords []string, location string) job.Scraper {
	return func() ([]job.Job, error) {
		var allJobs []job.Job

		// Build search URL - Dice uses /jobs/pages/ for pagination
		baseURL := "https://www.dice.com/jobs"
		queryParams := ""

		if len(keywords) > 0 {
			query := strings.Join(keywords, " ")
			queryParams += fmt.Sprintf("q=%s", strings.ReplaceAll(query, " ", "+"))
		}

		if location != "" {
			if queryParams != "" {
				queryParams += "&"
			}
			queryParams += fmt.Sprintf("location=%s", strings.ReplaceAll(location, " ", "+"))
		}

		// Try first few pages (be respectful with rate limiting)
		for page := 1; page <= 3; page++ {
			var url string
			if page == 1 {
				if queryParams != "" {
					url = fmt.Sprintf("%s?%s", baseURL, queryParams)
				} else {
					url = baseURL
				}
			} else {
				if queryParams != "" {
					url = fmt.Sprintf("%s/pages/%d?%s", baseURL, page, queryParams)
				} else {
					url = fmt.Sprintf("%s/pages/%d", baseURL, page)
				}
			}

			jobs, err := scrapeDicePage(url)
			if err != nil {
				// If we get an error, stop pagination
				break
			}

			if len(jobs) == 0 {
				// No more jobs found, stop pagination
				break
			}

			allJobs = append(allJobs, jobs...)

			// Rate limiting - be respectful to Dice.com
			time.Sleep(2 * time.Second)
		}

		return allJobs, nil
	}
}

func scrapeDicePage(url string) ([]job.Job, error) {
	scraper := BrowserScrape(url, func(page *rod.Page) ([]job.Job, error) {
		// Wait for page to stabilize
		page.MustWaitStable()

		// Wait a bit more for dynamic content to load
		time.Sleep(3 * time.Second)

		// Look for job cards - Dice uses various selectors over time
		jobSelectors := []string{
			"[data-cy='search-result-card']",
			".search-result-card",
			".card-search-result",
			"[data-testid='search-result']",
			".job-card",
		}

		var elements rod.Elements
		var err error

		// Try different selectors to find job cards
		for _, selector := range jobSelectors {
			elements, err = page.Elements(selector)
			if err == nil && len(elements) > 0 {
				break
			}
		}

		if err != nil || len(elements) == 0 {
			// No job cards found, try to check if there's a different layout
			return nil, fmt.Errorf("dice: no job cards found with any selector")
		}

		var jobs []job.Job

		for _, el := range elements {
			// Extract job title
			title, _ := el.Element("[data-cy='search-result-job-title'], .job-title, h3 a, h2 a, .card-title")
			if title == nil {
				continue
			}

			titleText, _ := title.Text()
			if titleText == "" {
				continue
			}

			// Extract company name
			company, _ := el.Element("[data-cy='search-result-company-name'], .company-name, .employer-name, [data-testid='company-name']")
			companyText := ""
			if company != nil {
				companyText, _ = company.Text()
			}

			// Extract location
			location, _ := el.Element("[data-cy='search-result-location'], .location, .job-location, [data-testid='location']")
			locationText := ""
			if location != nil {
				locationText, _ = location.Text()
			}

			// Extract job link
			var jobURL string
			if link, _ := el.Element("a"); link != nil {
				if href, _ := link.Attribute("href"); href != nil && *href != "" {
					jobURL = *href
					if !strings.HasPrefix(jobURL, "http") {
						jobURL = "https://www.dice.com" + jobURL
					}
				}
			}

			// Extract description/snippet
			description, _ := el.Element("[data-cy='search-result-description'], .job-description, .card-description, .job-summary")
			descriptionText := ""
			if description != nil {
				descriptionText, _ = description.Text()
			}

			// Extract salary if available
			salary, _ := el.Element("[data-cy='search-result-salary'], .salary, .pay-range, [data-testid='salary']")
			salaryText := ""
			if salary != nil {
				salaryText, _ = salary.Text()
			}

			// Extract posted date if available
			postedDate, _ := el.Element("[data-cy='search-result-date'], .posted-date, .date-posted, [data-testid='posted-date']")
			if postedDate != nil {
				_, _ = postedDate.Text() // Extract but don't use for now
			}

			// Create job object
			j := job.Job{
				ID:          idFromContent("dice", titleText+companyText+locationText),
				Title:       strings.TrimSpace(titleText),
				Company:     strings.TrimSpace(companyText),
				Location:    strings.TrimSpace(locationText),
				Description: strings.TrimSpace(descriptionText),
				URL:         jobURL,
				Source:      "dice",
				PostedDate:  time.Now(),
				Salary:      strings.TrimSpace(salaryText),
				Score:       50,
			}

			// Extract additional information from description
			if descriptionText != "" {
				j.Email = parse.ExtractFirstEmail(descriptionText)
				if salaryText == "" {
					j.Salary = parse.ExtractSalary(descriptionText)
				}
			}

			// Check for remote work indicators
			if strings.Contains(strings.ToLower(descriptionText+titleText+locationText), "remote") {
				j.JobType = "Remote"
			}

			jobs = append(jobs, j)
		}

		return jobs, nil
	})

	return scraper()
}
