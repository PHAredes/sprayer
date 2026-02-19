package job

import (
	"sort"
	"sprayer/internal/parse"
)

func Map(jobs []Job, f func(Job) Job) []Job {
	out := make([]Job, len(jobs))
	for i, j := range jobs {
		out[i] = f(j)
	}
	return out
}

func Select(jobs []Job, pred func(Job) bool) []Job {
	var out []Job
	for _, j := range jobs {
		if pred(j) {
			out = append(out, j)
		}
	}
	return out
}

func SortBy(less func(a, b Job) bool) Filter {
	return func(jobs []Job) []Job {
		sorted := make([]Job, len(jobs))
		copy(sorted, jobs)
		sort.Slice(sorted, func(i, j int) bool {
			return less(sorted[i], sorted[j])
		})
		return sorted
	}
}

var (
	ByScoreDesc = func(a, b Job) bool { return a.Score > b.Score }
	ByDateDesc  = func(a, b Job) bool { return a.PostedDate.After(b.PostedDate) }
	ByTitleAsc  = func(a, b Job) bool { return a.Title < b.Title }
)

func FlagTraps() Filter {
	return func(jobs []Job) []Job {
		return Map(jobs, func(j Job) Job {
			traps := parse.CheckForTraps(j.Description)
			if len(traps) > 0 {
				j.HasTraps = true
				j.Traps = traps
			}
			return j
		})
	}
}

func SanitizeDescriptions() Filter {
	return func(jobs []Job) []Job {
		return Map(jobs, func(j Job) Job {
			j.Description = parse.Sanitize(j.Description)
			return j
		})
	}
}

func Dedup() Filter {
	return func(jobs []Job) []Job {
		seen := make(map[string]bool, len(jobs))
		var out []Job
		for _, j := range jobs {
			if !seen[j.ID] {
				seen[j.ID] = true
				out = append(out, j)
			}
		}
		return out
	}
}
