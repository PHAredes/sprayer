package scraper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sprayer/internal/job"
	"sprayer/internal/parse"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// ExtractFn extracts jobs from a rod page. Higher-order: each site provides its own.
type ExtractFn func(page *rod.Page) ([]job.Job, error)

// BrowserScrape is a higher-order scraper: takes a URL and an extraction function,
// returns a Scraper. The browser is shared across calls.
func BrowserScrape(url string, extract ExtractFn) job.Scraper {
	return func() ([]job.Job, error) {
		l, err := launcher.New().Headless(true).Launch()
		if err != nil {
			return nil, fmt.Errorf("launch browser: %w", err)
		}

		browser := rod.New().ControlURL(l)
		if err := browser.Connect(); err != nil {
			return nil, fmt.Errorf("connect browser: %w", err)
		}
		defer browser.Close()

		page, err := browser.Page(proto.TargetCreateTarget{})
		if err != nil {
			return nil, fmt.Errorf("new page: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		page = page.Context(ctx)

		if err := page.Navigate(url); err != nil {
			return nil, fmt.Errorf("navigate to %s: %w", url, err)
		}

		if err := page.WaitLoad(); err != nil {
			return nil, fmt.Errorf("wait load: %w", err)
		}

		// Give JS some extra rendering time
		time.Sleep(2 * time.Second)

		return extract(page)
	}
}

// LinkedIn returns a browser-based scraper for LinkedIn job search.
func LinkedIn(keywords []string, location string) job.Scraper {
	query := strings.Join(keywords, " ")
	url := fmt.Sprintf("https://www.linkedin.com/jobs/search/?keywords=%s&location=%s&f_WT=2",
		strings.ReplaceAll(query, " ", "%20"),
		strings.ReplaceAll(location, " ", "%20"),
	)

	return BrowserScrape(url, func(page *rod.Page) ([]job.Job, error) {
		// Wait for job cards to load
		page.MustWaitStable()

		// Extract job cards from LinkedIn's public job search
		elements, err := page.Elements(".base-card")
		if err != nil {
			return nil, fmt.Errorf("find job cards: %w", err)
		}

		var jobs []job.Job
		for _, el := range elements {
			title, _ := el.Element(".base-search-card__title")
			company, _ := el.Element(".base-search-card__subtitle")
			loc, _ := el.Element(".job-search-card__location")
			link, _ := el.Element("a")

			if title == nil || company == nil {
				continue
			}

			titleText, _ := title.Text()
			companyText, _ := company.Text()
			locText := ""
			if loc != nil {
				locText, _ = loc.Text()
			}
			
			href := ""
			if link != nil {
				h, _ := link.Attribute("href")
				if h != nil {
					href = *h
				} else {
					// Fallback to property
					p, err := link.Property("href")
					if err == nil {
						href = p.String()
					}
				}
			}

			j := job.Job{
				ID:         idFromContent("li", titleText+companyText),
				Title:      strings.TrimSpace(titleText),
				Company:    strings.TrimSpace(companyText),
				Location:   strings.TrimSpace(locText),
				URL:        href,
				Source:     "linkedin",
				PostedDate: time.Now(),
				Score:      50,
			}

			// Try to get description from the detail
			if desc, _ := el.Text(); desc != "" {
				j.Description = strings.TrimSpace(desc)
				j.Email = parse.ExtractFirstEmail(desc)
				j.Salary = parse.ExtractSalary(desc)
			}

			jobs = append(jobs, j)
		}

		return jobs, nil
	})
}
