package job_test

import (
	"testing"
	"time"

	"sprayer/internal/job"
)

func now() time.Time { return time.Now() }

func TestPipeline_Composition(t *testing.T) {
	jobs := []job.Job{
		{ID: "1", Title: "Rust Dev", Company: "HOC", Score: 80, PostedDate: now()},
		{ID: "2", Title: "Go Dev", Company: "Google", Score: 90, PostedDate: now()},
		{ID: "3", Title: "Java Dev", Company: "Oracle", Score: 50, PostedDate: now()},
		{ID: "1", Title: "Rust Dev", Company: "HOC", Score: 80, PostedDate: now()}, // Duplicate ID=1
	}

	// Definition of pipeline: Dedup -> Filter(Rust) -> Sort(Score)
	pipeline := job.Pipe(
		job.Dedup(),
		job.ByKeywords([]string{"rust"}),
		job.SortBy(job.ByScoreDesc),
	)

	result := pipeline(jobs)

	if len(result) != 1 {
		t.Fatalf("Expected 1 job, got %d", len(result))
	}
	if result[0].Title != "Rust Dev" {
		t.Errorf("Expected Title 'Rust Dev', got '%s'", result[0].Title)
	}
}
