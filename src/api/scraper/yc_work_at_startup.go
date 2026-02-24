package scraper

import (
	"fmt"
	"strings"
	"time"

	"sprayer/src/api/job"
	"sprayer/src/api/parse"

	"github.com/go-rod/rod"
)

// YCWorkAtStartup returns a browser-based scraper for YC Work at a Startup.
// It navigates the React + Algolia search interface and extracts job listings.
func YCWorkAtStartup(keywords []string, location string) job.Scraper {
	url := "https://www.workatastartup.com/"

	return BrowserScrape(url, func(page *rod.Page) ([]job.Job, error) {
		// Wait for React app to load
		page.MustWaitStable()

		// Give extra time for Algolia search interface to initialize
		time.Sleep(4 * time.Second)

		// Extract job listings directly from the page
		return extractYCWorkAtStartupJobs(page)
	})
}

// extractYCWorkAtStartupJobs extracts job listings from the YC Work at a Startup page
func extractYCWorkAtStartupJobs(page *rod.Page) ([]job.Job, error) {
	// Look for job cards - YC uses various selectors for job listings
	jobSelectors := []string{
		"[data-testid='job-card']",
		".job-card",
		"[class*='job']",
		"[class*='Job']",
		".company-job",
		"[data-testid='company-job']",
	}

	var jobElements []*rod.Element
	var err error

	// Try different selectors to find job listings
	for _, selector := range jobSelectors {
		jobElements, err = page.Elements(selector)
		if err == nil && len(jobElements) > 0 {
			break
		}
	}

	if err != nil || len(jobElements) == 0 {
		// Fallback: look for any clickable job links
		jobElements, err = page.Elements("a[href*='/jobs/'], a[href*='/companies/']")
		if err != nil || len(jobElements) == 0 {
			return nil, fmt.Errorf("no job listings found")
		}
	}

	var jobs []job.Job

	for _, element := range jobElements {
		j, err := parseYCWorkAtStartupJob(element)
		if err != nil {
			continue // Skip jobs that can't be parsed
		}

		if j.Title != "" && j.Company != "" {
			jobs = append(jobs, j)
		}
	}

	return jobs, nil
}

// parseYCWorkAtStartupJob parses a single job element from YC Work at a Startup
func parseYCWorkAtStartupJob(element *rod.Element) (job.Job, error) {
	var j job.Job
	j.Source = "yc_work_at_startup"
	j.PostedDate = time.Now()
	j.Score = 80 // YC jobs are high quality

	// Extract job title
	titleSelectors := []string{
		"[data-testid='job-title']",
		".job-title",
		"h3", "h2", "h4",
		"[class*='title']",
		"[class*='Title']",
	}

	for _, selector := range titleSelectors {
		titleEl, err := element.Element(selector)
		if err == nil {
			title, err := titleEl.Text()
			if err == nil && title != "" {
				j.Title = strings.TrimSpace(title)
				break
			}
		}
	}

	// Extract company name
	companySelectors := []string{
		"[data-testid='company-name']",
		".company-name",
		"[class*='company']",
		"[class*='Company']",
	}

	for _, selector := range companySelectors {
		companyEl, err := element.Element(selector)
		if err == nil {
			company, err := companyEl.Text()
			if err == nil && company != "" {
				j.Company = strings.TrimSpace(company)
				break
			}
		}
	}

	// Extract location
	locationSelectors := []string{
		"[data-testid='job-location']",
		".job-location",
		"[class*='location']",
		"[class*='Location']",
	}

	for _, selector := range locationSelectors {
		locationEl, err := element.Element(selector)
		if err == nil {
			location, err := locationEl.Text()
			if err == nil && location != "" {
				j.Location = strings.TrimSpace(location)
				break
			}
		}
	}

	// Extract job URL
	link, err := element.Element("a")
	if err == nil {
		href, err := link.Attribute("href")
		if err == nil && href != nil {
			j.URL = *href
			if !strings.HasPrefix(j.URL, "http") {
				j.URL = "https://www.workatastartup.com" + j.URL
			}
		}
	}

	// Extract description from element text
	elementText, err := element.Text()
	if err == nil && elementText != "" {
		j.Description = strings.TrimSpace(elementText)
		j.Email = parse.ExtractFirstEmail(j.Description)
		j.Salary = parse.ExtractSalary(j.Description)

		// Extract YC batch information
		if strings.Contains(j.Description, "YC") {
			// Look for batch patterns like "YC S23", "YC W24", etc.
			if batch := extractYCBatch(j.Description); batch != "" {
				j.Description = fmt.Sprintf("YC Batch: %s\n%s", batch, j.Description)
			}
		}
	}

	// Extract job type
	jobTypeSelectors := []string{
		"[data-testid='job-type']",
		".job-type",
		"[class*='type']",
		"[class*='Type']",
	}

	for _, selector := range jobTypeSelectors {
		typeEl, err := element.Element(selector)
		if err == nil {
			jobType, err := typeEl.Text()
			if err == nil && jobType != "" {
				j.JobType = strings.TrimSpace(jobType)
				break
			}
		}
	}

	// Generate ID from content
	j.ID = idFromContent("yc_work_at_startup", j.Title+j.Company+j.Location)

	return j, nil
}

// extractYCBatch extracts YC batch information from text
func extractYCBatch(text string) string {
	// Look for patterns like "YC S23", "YC W24", "YCS23", "YCW24", etc.
	batchPatterns := []string{
		`YC\s*[SW]\d{2}`,
		`YCS\d{2}`,
		`YCW\d{2}`,
	}

	for _, pattern := range batchPatterns {
		// Simple pattern matching for YC batch
		if strings.Contains(text, pattern) {
			return pattern
		}
	}

	return ""
}
