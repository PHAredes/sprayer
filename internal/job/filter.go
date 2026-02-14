package job

import "strings"

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
