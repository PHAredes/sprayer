package scraper

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"sprayer/src/api/job"
	"sprayer/src/api/parse"
)

// AuthenticJobs scrapes the Authentic Jobs RSS feed.
func AuthenticJobs() job.Scraper {
	return func() ([]job.Job, error) {
		// Implement 3-second crawl delay to respect rate limiting
		time.Sleep(3 * time.Second)

		data, err := httpGet("https://authenticjobs.com/?feed=job_feed")
		if err != nil {
			return nil, fmt.Errorf("AuthenticJobs RSS: %w", err)
		}

		var feed authenticJobsFeed
		if err := xml.Unmarshal(data, &feed); err != nil {
			return nil, fmt.Errorf("AuthenticJobs parse: %w", err)
		}

		var jobs []job.Job
		for _, item := range feed.Channel.Items {
			desc := stripHTML(item.Description)
			content := stripHTML(item.ContentEncoded)
			fullDescription := desc + " " + content

			posted, _ := time.Parse(time.RFC1123Z, item.PubDate)
			if posted.IsZero() {
				posted, _ = time.Parse(time.RFC1123, item.PubDate)
			}
			if posted.IsZero() {
				posted = time.Now()
			}

			// Extract company, location, and job type from the raw XML content
			// Since namespace parsing isn't working, we'll extract from description and title
			company := extractCompanyFromAuthenticTitle(item.Title)
			if company == "" {
				company = extractCompanyFromURL(item.Link)
			}

			// Extract location from description using the parse package
			locations := parse.ExtractLocations(fullDescription)
			location := ""
			if len(locations) > 0 {
				location = strings.Join(locations, ", ")
			} else {
				// Fallback: look for common location patterns in description
				location = extractLocationFromText(fullDescription)
			}

			jobType := ""

			// Try to extract job type from description content
			lowerDesc := strings.ToLower(fullDescription)
			if strings.Contains(lowerDesc, "full-time") || strings.Contains(lowerDesc, "full time") {
				jobType = "Full-time"
			} else if strings.Contains(lowerDesc, "part-time") || strings.Contains(lowerDesc, "part time") {
				jobType = "Part-time"
			} else if strings.Contains(lowerDesc, "contract") {
				jobType = "Contract"
			} else if strings.Contains(lowerDesc, "freelance") {
				jobType = "Freelance"
			}

			j := job.Job{
				ID:          idFromContent("authenticjobs", item.Link+item.Title),
				Title:       item.Title,
				Company:     company,
				Location:    location,
				Description: fullDescription,
				URL:         item.Link,
				Source:      "authenticjobs",
				PostedDate:  posted,
				JobType:     jobType,
				Email:       parse.ExtractFirstEmail(fullDescription),
				Salary:      parse.ExtractSalary(fullDescription),
				Score:       50,
			}
			jobs = append(jobs, j)
		}
		return jobs, nil
	}
}

func extractCompanyFromAuthenticTitle(title string) string {
	// Many RSS job feeds use "Role at Company" or "Company - Role"
	if idx := strings.Index(title, " at "); idx > 0 {
		return strings.TrimSpace(title[idx+4:])
	}
	if parts := strings.SplitN(title, " - ", 2); len(parts) == 2 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}

func extractCompanyFromURL(url string) string {
	// Extract company from URL like: https://authenticjobs.com/job/35671/asana-engineering-manager-track-anything/
	// The URL structure is: /job/{job-id}/{company-name}-{job-title}/
	parts := strings.Split(url, "/")
	if len(parts) >= 5 {
		// The company name and job title are in the last part
		lastPart := parts[len(parts)-1]
		if lastPart == "" && len(parts) >= 6 {
			// Handle trailing slash
			lastPart = parts[len(parts)-2]
		}

		// Remove any query parameters
		lastPart = strings.Split(lastPart, "?")[0]

		// Common company names that we can identify
		commonCompanies := map[string]bool{
			"asana": true, "reddit": true, "twitch": true, "ramp": true,
			"openai": true, "anthropic": true, "google": true, "facebook": true,
			"amazon": true, "microsoft": true, "apple": true, "netflix": true,
		}

		// Split by dash and try to identify company name
		urlParts := strings.Split(lastPart, "-")
		for i, part := range urlParts {
			lowerPart := strings.ToLower(part)
			if commonCompanies[lowerPart] {
				// Found a known company name
				return strings.Title(part)
			}
			// If it's the first part and not a common job word, it might be the company
			if i == 0 && !isCommonJobWord(lowerPart) {
				return strings.Title(part)
			}
		}
	}
	return ""
}

func isCommonJobWord(word string) bool {
	commonWords := []string{"engineer", "developer", "manager", "designer", "scientist",
		"analyst", "specialist", "lead", "director", "coordinator", "assistant",
		"senior", "junior", "staff", "principal", "head", "vp", "vice", "president"}

	for _, common := range commonWords {
		if word == common {
			return true
		}
	}
	return false
}

func extractLocationFromText(text string) string {
	// Look for common location indicators
	lowerText := strings.ToLower(text)

	// Common city patterns
	cities := []string{"new york", "san francisco", "los angeles", "chicago", "boston",
		"seattle", "austin", "denver", "miami", "atlanta", "philadelphia", "washington",
		"london", "paris", "berlin", "amsterdam", "warsaw", "toronto", "vancouver"}

	for _, city := range cities {
		if strings.Contains(lowerText, city) {
			return strings.Title(city)
		}
	}

	// Look for "remote" or "hybrid"
	if strings.Contains(lowerText, "remote") {
		return "Remote"
	}
	if strings.Contains(lowerText, "hybrid") {
		return "Hybrid"
	}

	return ""
}

// authenticJobsFeed represents the Authentic Jobs RSS feed structure
type authenticJobsFeed struct {
	XMLName xml.Name             `xml:"rss"`
	Channel authenticJobsChannel `xml:"channel"`
}

type authenticJobsChannel struct {
	XMLName xml.Name            `xml:"channel"`
	Items   []authenticJobsItem `xml:"item"`
}

type authenticJobsItem struct {
	XMLName        xml.Name `xml:"item"`
	Title          string   `xml:"title"`
	Link           string   `xml:"link"`
	Description    string   `xml:"description"`
	PubDate        string   `xml:"pubDate"`
	ContentEncoded string   `xml:"content:encoded"`
	Company        string   `xml:"job_listing:company"`
	Location       string   `xml:"job_listing:location"`
	JobType        string   `xml:"job_listing:job_type"`
}
