package profile

import (
	"strings"

	"sprayer/internal/job"
)

// Scorer rates how well a job matches a profile. Composable.
type Scorer func(job.Job, Profile) float64

// CombineScorers averages multiple scorers (equal weight).
func CombineScorers(scorers ...Scorer) Scorer {
	return func(j job.Job, p Profile) float64 {
		if len(scorers) == 0 {
			return 0
		}
		var total float64
		for _, s := range scorers {
			total += s(j, p)
		}
		return total / float64(len(scorers))
	}
}

// KeywordScorer scores by keyword overlap with title + description.
func KeywordScorer(j job.Job, p Profile) float64 {
	if len(p.Keywords) == 0 {
		return 0
	}
	text := strings.ToLower(j.Title + " " + j.Description)
	var hits int
	for _, kw := range p.Keywords {
		if strings.Contains(text, strings.ToLower(kw)) {
			hits++
		}
	}
	return float64(hits) / float64(len(p.Keywords))
}

// LocationScorer scores 1.0 if locations match, 0.5 for remote, 0 otherwise.
func LocationScorer(j job.Job, p Profile) float64 {
	jLoc := strings.ToLower(j.Location)

	if strings.Contains(jLoc, "remote") && p.PreferRemote {
		return 1.0
	}

	for _, loc := range p.Locations {
		if strings.Contains(jLoc, strings.ToLower(loc)) {
			return 1.0
		}
	}

	if p.PreferRemote && jLoc == "" {
		return 0.3
	}

	return 0
}

// TitleScorer scores by keyword presence in the job title only.
func TitleScorer(j job.Job, p Profile) float64 {
	if len(p.Keywords) == 0 {
		return 0
	}
	title := strings.ToLower(j.Title)
	var hits int
	for _, kw := range p.Keywords {
		if strings.Contains(title, strings.ToLower(kw)) {
			hits++
		}
	}
	return float64(hits) / float64(len(p.Keywords))
}

// BestProfile returns the profile that best matches a job, using the provided scorer.
func BestProfile(j job.Job, profiles []Profile, scorer Scorer) (Profile, float64) {
	var best Profile
	var bestScore float64

	for _, p := range profiles {
		score := scorer(j, p)
		if score > bestScore {
			bestScore = score
			best = p
		}
	}

	return best, bestScore
}

// DefaultScorer is the standard scorer combining keyword, title, and location.
var DefaultScorer = CombineScorers(KeywordScorer, TitleScorer, LocationScorer)
