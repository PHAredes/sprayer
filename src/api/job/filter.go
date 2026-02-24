package job

import (
	"strings"
	"time"
)

// Filter transforms a job list. Chainable via Pipe().
type Filter func([]Job) []Job

// Pipe composes filters left-to-right: Pipe(f, g)(jobs) == g(f(jobs)).
func Pipe(filters ...Filter) Filter {
	return func(jobs []Job) []Job {
		for _, f := range filters {
			jobs = f(jobs)
		}
		return jobs
	}
}

// ByKeywords returns jobs matching any keyword in title or description.
func ByKeywords(keywords []string) Filter {
	return func(jobs []Job) []Job {
		if len(keywords) == 0 {
			return jobs
		}
		var out []Job
		for _, j := range jobs {
			lower := strings.ToLower(j.Title + " " + j.Description)
			for _, kw := range keywords {
				if strings.Contains(lower, strings.ToLower(strings.TrimSpace(kw))) {
					out = append(out, j)
					break
				}
			}
		}
		return out
	}
}

// ExcludeKeywords filters out jobs containing any of the specified keywords
func ExcludeKeywords(keywords []string) Filter {
	return func(jobs []Job) []Job {
		if len(keywords) == 0 {
			return jobs
		}
		var out []Job
		for _, j := range jobs {
			lower := strings.ToLower(j.Title + " " + j.Description)
			excluded := false
			for _, kw := range keywords {
				if strings.Contains(lower, strings.ToLower(strings.TrimSpace(kw))) {
					excluded = true
					break
				}
			}
			if !excluded {
				out = append(out, j)
			}
		}
		return out
	}
}

// ByMinScore returns jobs with score >= min.
func ByMinScore(min int) Filter {
	return func(jobs []Job) []Job {
		if min <= 0 {
			return jobs
		}
		var out []Job
		for _, j := range jobs {
			if j.Score >= min {
				out = append(out, j)
			}
		}
		return out
	}
}

// ByScoreRange returns jobs within the specified score range
func ByScoreRange(min, max int) Filter {
	return func(jobs []Job) []Job {
		if min <= 0 && max >= 100 {
			return jobs
		}
		var out []Job
		for _, j := range jobs {
			if j.Score >= min && j.Score <= max {
				out = append(out, j)
			}
		}
		return out
	}
}

// ByLocation returns jobs matching a location substring.
func ByLocation(loc string) Filter {
	return func(jobs []Job) []Job {
		if loc == "" {
			return jobs
		}
		lower := strings.ToLower(loc)
		var out []Job
		for _, j := range jobs {
			if strings.Contains(strings.ToLower(j.Location), lower) {
				out = append(out, j)
			}
		}
		return out
	}
}

// ByLocations returns jobs matching any of the specified locations
func ByLocations(locations []string) Filter {
	return func(jobs []Job) []Job {
		if len(locations) == 0 {
			return jobs
		}
		var out []Job
		for _, j := range jobs {
			for _, loc := range locations {
				if strings.Contains(strings.ToLower(j.Location), strings.ToLower(strings.TrimSpace(loc))) {
					out = append(out, j)
					break
				}
			}
		}
		return out
	}
}

// ByCompany returns jobs matching a company substring.
func ByCompany(company string) Filter {
	return func(jobs []Job) []Job {
		if company == "" {
			return jobs
		}
		lower := strings.ToLower(company)
		var out []Job
		for _, j := range jobs {
			if strings.Contains(strings.ToLower(j.Company), lower) {
				out = append(out, j)
			}
		}
		return out
	}
}

// ByCompanies returns jobs matching any of the specified companies
func ByCompanies(companies []string) Filter {
	return func(jobs []Job) []Job {
		if len(companies) == 0 {
			return jobs
		}
		var out []Job
		for _, j := range jobs {
			for _, company := range companies {
				if strings.Contains(strings.ToLower(j.Company), strings.ToLower(strings.TrimSpace(company))) {
					out = append(out, j)
					break
				}
			}
		}
		return out
	}
}

// ExcludeCompanies filters out jobs from specified companies
func ExcludeCompanies(companies []string) Filter {
	return func(jobs []Job) []Job {
		if len(companies) == 0 {
			return jobs
		}
		var out []Job
		for _, j := range jobs {
			excluded := false
			for _, company := range companies {
				if strings.Contains(strings.ToLower(j.Company), strings.ToLower(strings.TrimSpace(company))) {
					excluded = true
					break
				}
			}
			if !excluded {
				out = append(out, j)
			}
		}
		return out
	}
}

// HasEmail returns jobs that have an email field set.
func HasEmail() Filter {
	return func(jobs []Job) []Job {
		var out []Job
		for _, j := range jobs {
			if j.Email != "" {
				out = append(out, j)
			}
		}
		return out
	}
}

// ExcludeTraps filters out jobs that have traps
func ExcludeTraps() Filter {
	return func(jobs []Job) []Job {
		var out []Job
		for _, j := range jobs {
			if !j.HasTraps {
				out = append(out, j)
			}
		}
		return out
	}
}

// RemotePreferred prioritizes remote jobs
func RemotePreferred() Filter {
	return func(jobs []Job) []Job {
		var remote []Job
		var onsite []Job

		for _, j := range jobs {
			if strings.Contains(strings.ToLower(j.Location), "remote") {
				remote = append(remote, j)
			} else {
				onsite = append(onsite, j)
			}
		}

		// Return remote jobs first, then onsite
		return append(remote, onsite...)
	}
}

// BySeniorityLevels returns jobs matching specified seniority levels
func BySeniorityLevels(levels []string) Filter {
	return func(jobs []Job) []Job {
		if len(levels) == 0 {
			return jobs
		}
		var out []Job
		for _, j := range jobs {
			titleLower := strings.ToLower(j.Title)
			for _, level := range levels {
				if strings.Contains(titleLower, strings.ToLower(level)) {
					out = append(out, j)
					break
				}
			}
		}
		return out
	}
}

// ByTechnologies returns jobs mentioning specified technologies
func ByTechnologies(techs []string) Filter {
	return func(jobs []Job) []Job {
		if len(techs) == 0 {
			return jobs
		}
		var out []Job
		for _, j := range jobs {
			contentLower := strings.ToLower(j.Title + " " + j.Description)
			for _, tech := range techs {
				if strings.Contains(contentLower, strings.ToLower(strings.TrimSpace(tech))) {
					out = append(out, j)
					break
				}
			}
		}
		return out
	}
}

// ExcludeTechnologies filters out jobs mentioning specified technologies
func ExcludeTechnologies(techs []string) Filter {
	return func(jobs []Job) []Job {
		if len(techs) == 0 {
			return jobs
		}
		var out []Job
		for _, j := range jobs {
			contentLower := strings.ToLower(j.Title + " " + j.Description)
			excluded := false
			for _, tech := range techs {
				if strings.Contains(contentLower, strings.ToLower(strings.TrimSpace(tech))) {
					excluded = true
					break
				}
			}
			if !excluded {
				out = append(out, j)
			}
		}
		return out
	}
}

// PostedAfter filters jobs posted after the specified time
func PostedAfter(after time.Time) Filter {
	return func(jobs []Job) []Job {
		var out []Job
		for _, j := range jobs {
			if j.PostedDate.After(after) {
				out = append(out, j)
			}
		}
		return out
	}
}

// PostedBefore filters jobs posted before the specified time
func PostedBefore(before time.Time) Filter {
	return func(jobs []Job) []Job {
		var out []Job
		for _, j := range jobs {
			if j.PostedDate.Before(before) {
				out = append(out, j)
			}
		}
		return out
	}
}

// CVMatcher interface for CV-based job matching
type CVMatcher interface {
	ScoreJob(*Job) int
}

// ByCVMatch returns jobs that match the provided CV
func ByCVMatch(cv CVMatcher, minScore int) Filter {
	return func(jobs []Job) []Job {
		var filtered []Job
		for _, j := range jobs {
			score := cv.ScoreJob(&j)
			if score >= minScore {
				j.Score += score // Add CV score to existing score
				filtered = append(filtered, j)
			}
		}
		return filtered
	}
}
