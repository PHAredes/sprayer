package scraper

import (
	"context"
	"fmt"
	"sync"
	"time"

	"sprayer/internal/job"
	"sprayer/internal/profile"
)

// IncrementalScraper provides streaming job results with progress tracking
type IncrementalScraper struct {
	ctx           context.Context
	cancel        context.CancelFunc
	profile       profile.Profile
	results       chan job.Job
	errors        chan error
	progress      chan ScraperProgress
	totalJobs     int
	processedJobs int
	mu            sync.RWMutex
}

// Done returns a channel that's closed when the scraper is done
func (is *IncrementalScraper) Done() <-chan struct{} {
	return is.ctx.Done()
}

// ScraperProgress represents scraping progress information
type ScraperProgress struct {
	Source        string
	JobsFound     int
	TotalSources  int
	CurrentSource int
	ElapsedTime   time.Duration
	Status        string
}

// ScraperSource represents a scraper source with metadata
type ScraperSource struct {
	name string
	fn   ScraperFunc
}

// ScraperFunc represents a scraping function with context support
type ScraperFunc func(ctx context.Context, keywords []string, location string) ([]job.Job, error)

// NewIncrementalScraper creates a new incremental scraper
func NewIncrementalScraper(ctx context.Context, prof profile.Profile) *IncrementalScraper {
	ctx, cancel := context.WithCancel(ctx)

	return &IncrementalScraper{
		ctx:      ctx,
		cancel:   cancel,
		profile:  prof,
		results:  make(chan job.Job, 100),
		errors:   make(chan error, 10),
		progress: make(chan ScraperProgress, 10),
	}
}

// Start begins incremental scraping
func (is *IncrementalScraper) Start() {
	go is.runScraping()
}

// Stop cancels scraping
func (is *IncrementalScraper) Stop() {
	is.cancel()
}

// Results returns the job results channel
func (is *IncrementalScraper) Results() <-chan job.Job {
	return is.results
}

// Errors returns the error channel
func (is *IncrementalScraper) Errors() <-chan error {
	return is.errors
}

// Progress returns the progress channel
func (is *IncrementalScraper) Progress() <-chan ScraperProgress {
	return is.progress
}

// GetProgress returns current progress information
func (is *IncrementalScraper) GetProgress() (int, int, time.Duration) {
	is.mu.RLock()
	defer is.mu.RUnlock()
	return is.processedJobs, is.totalJobs, time.Since(time.Now())
}

func (is *IncrementalScraper) runScraping() {
	defer close(is.results)
	defer close(is.errors)
	defer close(is.progress)

	startTime := time.Now()
	sources := is.getScraperSources()

	is.mu.Lock()
	is.totalJobs = len(sources)
	is.mu.Unlock()

	// Create base keywords from profile
	keywords := is.profile.Keywords
	if len(keywords) == 0 {
		keywords = []string{"golang", "rust", "remote"}
	}

	location := ""
	if is.profile.PreferRemote {
		location = "Remote"
	}

	// Process sources incrementally
	for i, source := range sources {
		select {
		case <-is.ctx.Done():
			is.progress <- ScraperProgress{
				Status:      "Cancelled",
				ElapsedTime: time.Since(startTime),
			}
			return
		default:
		}

		sourceName := source.name
		is.sendProgress(sourceName, 0, len(sources), i+1, time.Since(startTime), "Scraping")

		// Run scraper with timeout
		ctx, cancel := context.WithTimeout(is.ctx, 30*time.Second)
		jobs, err := source.fn(ctx, keywords, location)
		cancel()

		if err != nil {
			is.errors <- fmt.Errorf("error scraping %s: %w", sourceName, err)
			continue
		}

		// Apply profile scoring and filtering incrementally
		filteredJobs := is.processJobsIncrementally(jobs)

		is.mu.Lock()
		is.processedJobs++
		is.mu.Unlock()

		is.sendProgress(sourceName, len(filteredJobs), len(sources), i+1, time.Since(startTime), "Complete")

		// Send results as they're processed
		for _, job := range filteredJobs {
			select {
			case is.results <- job:
			case <-is.ctx.Done():
				return
			}
		}
	}

	is.progress <- ScraperProgress{
		Status:      "Finished",
		ElapsedTime: time.Since(startTime),
	}
}

func (is *IncrementalScraper) processJobsIncrementally(jobs []job.Job) []job.Job {
	var filteredJobs []job.Job

	for _, j := range jobs {
		// Apply profile scoring
		j.Score = is.profile.CalculateJobScore(&j)

		// Apply basic filters
		if j.Score < is.profile.MinScore || j.Score > is.profile.MaxScore {
			continue
		}

		if is.profile.ExcludeTraps && j.HasTraps {
			continue
		}

		filteredJobs = append(filteredJobs, j)
	}

	// Apply profile filter pipeline
	filters := is.profile.GenerateFilters()
	return job.Pipe(filters...)(filteredJobs)
}

func (is *IncrementalScraper) sendProgress(sourceName string, jobsFound, totalSources, currentSource int, elapsed time.Duration, status string) {
	select {
	case is.progress <- ScraperProgress{
		Source:        sourceName,
		JobsFound:     jobsFound,
		TotalSources:  totalSources,
		CurrentSource: currentSource,
		ElapsedTime:   elapsed,
		Status:        status,
	}:
	case <-is.ctx.Done():
	}
}

func (is *IncrementalScraper) getScraperSources() []ScraperSource {
	return []ScraperSource{
		{name: "Hacker News", fn: func(ctx context.Context, keywords []string, location string) ([]job.Job, error) {
			return HN()()
		}},
		{name: "RemoteOK", fn: func(ctx context.Context, keywords []string, location string) ([]job.Job, error) {
			return RemoteOK()()
		}},
		{name: "Greenhouse", fn: func(ctx context.Context, keywords []string, location string) ([]job.Job, error) {
			return Greenhouse(DefaultGreenhouseBoards)()
		}},
		{name: "We Work Remotely", fn: func(ctx context.Context, keywords []string, location string) ([]job.Job, error) {
			return WeWorkRemotely()()
		}},
		{name: "Arbeitnow", fn: func(ctx context.Context, keywords []string, location string) ([]job.Job, error) {
			return Arbeitnow()()
		}},
		{name: "Jobicy", fn: func(ctx context.Context, keywords []string, location string) ([]job.Job, error) {
			return Jobicy()()
		}},
		{name: "RSS Feeds", fn: func(ctx context.Context, keywords []string, location string) ([]job.Job, error) {
			// RSS feeds need to be handled differently - return empty for now
			return []job.Job{}, nil
		}},
		{name: "LinkedIn", fn: func(ctx context.Context, keywords []string, location string) ([]job.Job, error) {
			return LinkedIn(keywords, location)()
		}},
		{name: "Indeed", fn: func(ctx context.Context, keywords []string, location string) ([]job.Job, error) {
			return Indeed(keywords[0], location)()
		}},
		{name: "Glassdoor", fn: func(ctx context.Context, keywords []string, location string) ([]job.Job, error) {
			return Glassdoor(keywords[0])()
		}},
	}
}
