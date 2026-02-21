package scraper

import (
	"sprayer/internal/job"
)

// All returns a merged scraper that hits every source.
// API-based scrapers run first (fast), browser-based scrapers follow.
func All(keywords []string, location string) job.Scraper {
	query := ""
	if len(keywords) > 0 {
		query = keywords[0]
	}

	// API-based (fast, reliable)
	api := []job.Scraper{
		HN(),
		RemoteOK(),
		Greenhouse(DefaultGreenhouseBoards),
		WeWorkRemotely(),
		Arbeitnow(),
		Jobicy(),
	}

	// Add RSS feeds
	api = append(api, CommonRSSFeeds()...)

	// Browser-based (slower, JS-rendered)
	browser := []job.Scraper{
		LinkedIn(keywords, location),
		Indeed(query, location),
		Glassdoor(query),
	}

	// Merge all: API first, then browser
	all := append(api, browser...)
	return job.Merge(all...)
}

// APIOnly returns a merged scraper with only API-based sources (no browser needed).
func APIOnly() job.Scraper {
	api := []job.Scraper{
		HN(),
		RemoteOK(),
		Greenhouse(DefaultGreenhouseBoards),
		WeWorkRemotely(),
		Arbeitnow(),
		Jobicy(),
	}
	api = append(api, CommonRSSFeeds()...)
	return job.Merge(api...)
}
