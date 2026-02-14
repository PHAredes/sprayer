package scraper

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"sprayer/internal/job"
	"sprayer/internal/parse"
)

// RSS creates a scraper from any RSS/Atom job board feed.
// Higher-order: takes a source name and URL, returns a Scraper.
func RSS(source, feedURL string) job.Scraper {
	return func() ([]job.Job, error) {
		data, err := httpGet(feedURL)
		if err != nil {
			return nil, fmt.Errorf("RSS %s: %w", source, err)
		}

		var feed rssFeed
		if err := xml.Unmarshal(data, &feed); err != nil {
			return nil, fmt.Errorf("RSS %s parse: %w", source, err)
		}

		var jobs []job.Job
		for _, item := range feed.Channel.Items {
			desc := stripHTML(item.Description)
			posted, _ := time.Parse(time.RFC1123Z, item.PubDate)
			if posted.IsZero() {
				posted, _ = time.Parse(time.RFC1123, item.PubDate)
			}
			if posted.IsZero() {
				posted = time.Now()
			}

			j := job.Job{
				ID:          idFromContent(source, item.Link+item.Title),
				Title:       item.Title,
				Company:     extractCompanyFromTitle(item.Title),
				Location:    strings.Join(parse.ExtractLocations(desc), ", "),
				Description: desc,
				URL:         item.Link,
				Source:      source,
				PostedDate:  posted,
				Email:       parse.ExtractFirstEmail(desc),
				Salary:      parse.ExtractSalary(desc),
				Score:       50,
			}
			jobs = append(jobs, j)
		}
		return jobs, nil
	}
}

type rssFeed struct {
	Channel struct {
		Items []rssItem `xml:"item"`
	} `xml:"channel"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func extractCompanyFromTitle(title string) string {
	// Many RSS job feeds use "Role at Company" or "Company - Role"
	if idx := strings.Index(title, " at "); idx > 0 {
		return strings.TrimSpace(title[idx+4:])
	}
	if parts := strings.SplitN(title, " - ", 2); len(parts) == 2 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}

// CommonRSSFeeds returns scrapers for well-known RSS job feeds.
func CommonRSSFeeds() []job.Scraper {
	feeds := []struct {
		name string
		url  string
	}{
		{"crypto-jobs", "https://crypto.jobs/feed"},
		{"nodejs-jobs", "https://nodesk.co/remote-jobs/rss/"},
		{"golang-cafe", "https://golang.cafe/Ede/rss.xml"},
		{"rustjobs", "https://rustjobs.dev/feed.xml"},
		{"functional-works", "https://functional.works-hub.com/feed"},
	}

	var scrapers []job.Scraper
	for _, f := range feeds {
		scrapers = append(scrapers, RSS(f.name, f.url))
	}
	return scrapers
}
