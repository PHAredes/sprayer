package profile_test

import (
	"testing"

	"sprayer/internal/job"
	"sprayer/internal/profile"
)

func TestScorer_Property(t *testing.T) {
	// Property: Score is always between 0 and 1 (or normalized range)
	// Actually our score functions allow > 1.0 depending on implementation?
	// Let's check: KeywordScorer returns fraction 0..1?
	// "float64(hits) / float64(len(p.Keywords))" -> yes, 0 to 1.
	// Unless keywords is empty, then 0.

	p := profile.Profile{Keywords: []string{"go", "rust"}}
	j := job.Job{Title: "Go Developer", Description: "We use Rust"}

	score := profile.KeywordScorer(j, p)
	if score < 0 || score > 1.0 {
		t.Errorf("Score out of range [0,1]: %f", score)
	}
}

func TestBestProfile_Selection(t *testing.T) {
	p1 := profile.Profile{ID: "p1", Keywords: []string{"java"}}
	p2 := profile.Profile{ID: "p2", Keywords: []string{"rust"}}

	profiles := []profile.Profile{p1, p2}
	j := job.Job{Title: "Rust Engineer", Description: "Rust stuff"}

	best, score := profile.BestProfile(j, profiles, profile.DefaultScorer)

	if best.ID != "p2" {
		t.Errorf("Expected profile p2, got %s", best.ID)
	}
	if score <= 0 {
		t.Errorf("Expected positive score, got %f", score)
	}
}
