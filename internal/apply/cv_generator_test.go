package apply

import (
	"testing"

	"sprayer/internal/profile"
)

func TestCVGenerator_Cache(t *testing.T) {
	gen := NewCVGenerator(nil)

	gen.ClearCache()
	if len(gen.cache) != 0 {
		t.Errorf("expected empty cache after clear")
	}
}

func TestCVGenerator_Available(t *testing.T) {
	gen := NewCVGenerator(nil)
	if gen.Available() {
		t.Errorf("expected unavailable with nil client")
	}
}

func TestSaveCustomCV(t *testing.T) {
	tmpDir := t.TempDir()

	content := "Test CV Content"
	jobID := "test-job-123"

	path, err := SaveCustomCV(content, jobID, tmpDir)
	if err != nil {
		t.Fatalf("SaveCustomCV failed: %v", err)
	}

	if path == "" {
		t.Errorf("expected non-empty path")
	}
}

func TestFormatExperience(t *testing.T) {
	experiences := []profile.Experience{
		{
			Title:       "Senior Developer",
			Company:     "Tech Corp",
			Duration:    "2020-2023",
			Description: "Built microservices",
		},
	}

	result := formatExperience(experiences)

	if result == "" {
		t.Errorf("expected non-empty result")
	}
}

func TestFormatEducation(t *testing.T) {
	education := []profile.Education{
		{
			Degree:      "BS",
			Field:       "Computer Science",
			Institution: "University",
			Year:        "2018",
		},
	}

	result := formatEducation(education)

	if result == "" {
		t.Errorf("expected non-empty result")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a very long string", 10, "this is a ..."},
		{"exact", 5, "exact"},
	}

	for _, tt := range tests {
		result := truncate(tt.input, tt.max)
		if result != tt.expected {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, result, tt.expected)
		}
	}
}
