package scraper

import (
	"fmt"
	"time"

	"sprayer/src/api/job"
)

// RemoteCoBrowser scrapes Remote.co using browser automation as a fallback.
// This is useful if Remote.co doesn't have public APIs or RSS feeds.
//
// NOTE: As of 2026-02-20, Remote.co appears to be inaccessible from this environment
// due to potential geographic/IP-based restrictions. The browser automation
// approach may work if the accessibility issues are resolved.
func RemoteCoBrowser() job.Scraper {
	return func() ([]job.Job, error) {
		// Test connectivity first
		if err := TestRemoteCoConnectivity(); err != nil {
			return nil, fmt.Errorf("remote.co browser scraper: accessibility test failed: %w", err)
		}

		// Potential Remote.co pages to scrape
		urls := []string{
			"https://remote.co/remote-jobs/",
			"https://remote.co/remote-jobs/developer/",
			"https://remote.co/remote-jobs/design/",
			"https://remote.co/remote-jobs/marketing/",
			"https://remote.co/remote-jobs/customer-service/",
		}

		var allJobs []job.Job

		for _, url := range urls {
			jobs, err := scrapeRemoteCoPage(url)
			if err != nil {
				continue // Try next URL
			}
			allJobs = append(allJobs, jobs...)
			time.Sleep(1 * time.Second) // Be respectful
		}

		return allJobs, nil
	}
}

func scrapeRemoteCoPage(url string) ([]job.Job, error) {
	// This would use the browser automation from browser.go
	// For now, return a placeholder implementation

	// Expected structure based on typical job board layouts:
	// - Job listings in cards or table rows
	// - Each job has: title, company, location, description snippet
	// - Link to full job details
	// - Posted date

	return nil, fmt.Errorf("browser scraping not yet implemented for Remote.co")
}

// Expected Remote.co job structure (based on common patterns)
type remoteCoJobListing struct {
	Title       string
	Company     string
	Location    string
	Description string
	URL         string
	PostedDate  string
	JobType     string
	Salary      string
	Tags        []string
}

// CSS selectors that might work for Remote.co (need to be verified)
var remoteCoSelectors = struct {
	JobContainer string
	Title        string
	Company      string
	Location     string
	Description  string
	URL          string
	PostedDate   string
	Salary       string
	Tags         string
}{
	JobContainer: ".job-listing, .job-card, .job-item, tr.job", // Common job container selectors
	Title:        ".job-title, .title, h2 a, h3 a",             // Job title selectors
	Company:      ".company, .company-name, .employer",         // Company name selectors
	Location:     ".location, .job-location",                   // Location selectors
	Description:  ".description, .job-summary, .excerpt",       // Description selectors
	URL:          "a.job-link, a[href*='job'], h2 a, h3 a",     // Job link selectors
	PostedDate:   ".date, .posted, .time",                      // Date selectors
	Salary:       ".salary, .pay, .compensation",               // Salary selectors
	Tags:         ".tags, .categories, .badge",                 // Tag selectors
}

// Common job categories on Remote.co (estimated)
func getRemoteCoCategories() []string {
	return []string{
		"developer",        // Programming/Development
		"design",           // Design/UI/UX
		"marketing",        // Marketing
		"sales",            // Sales
		"customer-service", // Customer Support
		"writing",          // Content/Writing
		"product",          // Product Management
		"business",         // Business/Operations
		"finance",          // Finance/Accounting
		"legal",            // Legal
		"data",             // Data Analysis/Science
		"qa",               // Quality Assurance
		"devops",           // DevOps/SysAdmin
		"all-others",       // Other categories
	}
}
