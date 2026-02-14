package job

// Scraper fetches jobs from a source. Composable: combine with Merge().
type Scraper func() ([]Job, error)

// Merge combines multiple scrapers into one. Errors are collected, not fatal.
func Merge(scrapers ...Scraper) Scraper {
	return func() ([]Job, error) {
		type result struct {
			jobs []Job
			err  error
		}
		ch := make(chan result, len(scrapers))
		
		for _, s := range scrapers {
			go func(s Scraper) {
				jobs, err := s()
				ch <- result{jobs, err}
			}(s)
		}

		var all []Job
		var lastErr error
		for i := 0; i < len(scrapers); i++ {
			res := <-ch
			if res.err != nil {
				lastErr = res.err
				// Log error but continue
				continue
			}
			all = append(all, res.jobs...)
		}
		
		return all, lastErr
	}
}
