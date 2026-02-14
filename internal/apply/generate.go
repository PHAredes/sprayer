package apply

import (
	"fmt"
	"strings"

	"sprayer/internal/job"
	"sprayer/internal/llm"
	"sprayer/internal/parse"
	"sprayer/internal/profile"
)

// GenerateEmail uses syntactic parsing + LLM to produce a personalized application email.
// Returns subject and body.
func GenerateEmail(j job.Job, p profile.Profile, client *llm.Client, promptName string) (string, string, error) {
	// 1. Extract context via syntactic parsing
	email := j.Email
	if email == "" {
		email = parse.ExtractFirstEmail(j.Description)
	}

	location := j.Location
	if location == "" {
		locs := parse.ExtractLocations(j.Description)
		if len(locs) > 0 {
			location = locs[0]
		}
	}

	// 2. Build prompt variables
	vars := map[string]string{
		"job_title":       j.Title,
		"company":         j.Company,
		"location":        location,
		"applicant_name":  p.Name,
		"skills":          strings.Join(p.Keywords, ", "),
		"job_description": truncate(j.Description, 2000),
		"applied_date":    j.AppliedDate.Format("2006-01-02"),
	}

	// 3. Load and interpolate prompt
	prompt, err := llm.LoadPrompt(promptName, vars)
	if err != nil {
		return "", "", fmt.Errorf("load prompt %q: %w", promptName, err)
	}

	// 4. Generate via LLM
	body, err := client.Complete(
		"You are a professional job application assistant. Be concise and natural.",
		prompt,
	)
	if err != nil {
		return "", "", fmt.Errorf("LLM generation: %w", err)
	}

	subject := fmt.Sprintf("Application for %s — %s", j.Title, p.Name)

	return subject, body, nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
