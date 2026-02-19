package scraper

import (
	"sprayer/internal/job"
	"sprayer/internal/profile"
)

// ProfileBasedScraper creates a scraper based on profile preferences
func ProfileBasedScraper(profile profile.Profile) job.Scraper {
	// Build keywords from profile
	keywords := profile.Keywords
	if len(keywords) == 0 {
		keywords = []string{"golang", "rust", "remote"} // Default fallback
	}

	// Determine location for scraping
	location := ""
	if profile.PreferRemote {
		location = "Remote"
	} else if len(profile.Locations) > 0 {
		location = profile.Locations[0] // Use first preferred location
	}

	// Create base scraper
	baseScraper := All(keywords, location)

	// Apply profile-based post-processing
	return func() ([]job.Job, error) {
		jobs, err := baseScraper()
		if err != nil {
			return nil, err
		}

		// Apply profile scoring and filtering
		scoredJobs := make([]job.Job, len(jobs))
		for i, job := range jobs {
			// Calculate custom score based on profile
			job.Score = profile.CalculateJobScore(&job)
			scoredJobs[i] = job
		}

		// Apply profile filters
		filters := profile.GenerateFilters()
		filteredJobs := job.Pipe(filters...)(scoredJobs)

		return filteredJobs, nil
	}
}

// FastProfileScraper creates a fast scraper using only API sources with profile preferences
func FastProfileScraper(profile profile.Profile) job.Scraper {
	keywords := profile.Keywords
	if len(keywords) == 0 {
		keywords = []string{"golang", "rust", "remote"}
	}

	// Use only fast API sources
	baseScraper := APIOnly()

	return func() ([]job.Job, error) {
		jobs, err := baseScraper()
		if err != nil {
			return nil, err
		}

		// Filter by keywords first (since API sources might not respect keywords)
		keywordFilter := job.ByKeywords(keywords)
		keywordJobs := keywordFilter(jobs)

		// Apply profile scoring and filtering
		scoredJobs := make([]job.Job, len(keywordJobs))
		for i, job := range keywordJobs {
			job.Score = profile.CalculateJobScore(&job)
			scoredJobs[i] = job
		}

		// Apply profile filters
		filters := profile.GenerateFilters()
		filteredJobs := job.Pipe(filters...)(scoredJobs)

		return filteredJobs, nil
	}
}

// SmartScraper creates an intelligent scraper that adapts to profile preferences
func SmartScraper(profile profile.Profile, fastMode bool) job.Scraper {
	if fastMode {
		return FastProfileScraper(profile)
	}
	return ProfileBasedScraper(profile)
}
