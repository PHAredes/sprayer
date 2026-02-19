package apply

import (
	"fmt"
	"strings"
	"sync"

	"sprayer/internal/job"
	"sprayer/internal/llm"
	"sprayer/internal/parse"
	"sprayer/internal/profile"
)

var (
	coverLetterCache = make(map[string]string)
	cacheMutex       sync.RWMutex
)

func GetCachedCoverLetter(jobID string) (string, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	content, ok := coverLetterCache[jobID]
	return content, ok
}

func CacheCoverLetter(jobID, content string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	coverLetterCache[jobID] = content
}

func ClearCoverLetterCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	coverLetterCache = make(map[string]string)
}

func GenerateEmail(j job.Job, p profile.Profile, client *llm.Client, promptName string) (string, string, error) {
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

	vars := map[string]string{
		"job_title":       j.Title,
		"company":         j.Company,
		"location":        location,
		"applicant_name":  p.Name,
		"skills":          strings.Join(p.Keywords, ", "),
		"job_description": truncate(parse.Sanitize(j.Description), 2000),
		"applied_date":    j.AppliedDate.Format("2006-01-02"),
	}

	prompt, err := llm.LoadPrompt(promptName, vars)
	if err != nil {
		return "", "", fmt.Errorf("load prompt %q: %w", promptName, err)
	}

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

func GenerateCoverLetter(j *job.Job, p *profile.Profile, client *llm.Client) (string, error) {
	if j == nil {
		return "", fmt.Errorf("job cannot be nil")
	}
	if p == nil {
		return "", fmt.Errorf("profile cannot be nil")
	}

	if cached, ok := GetCachedCoverLetter(j.ID); ok {
		return cached, nil
	}

	location := j.Location
	if location == "" {
		locs := parse.ExtractLocations(j.Description)
		if len(locs) > 0 {
			location = locs[0]
		}
	}

	sanitizedDesc := parse.Sanitize(j.Description)
	requirements := extractJobRequirements(sanitizedDesc)
	matchedTech := findMatchingTechnologies(j.Description, p)

	vars := map[string]string{
		"job_title":       j.Title,
		"company":         j.Company,
		"location":        location,
		"applicant_name":  p.Name,
		"job_description": truncate(sanitizedDesc, 2500),
		"requirements":    requirements,
		"matched_tech":    matchedTech,
	}

	vars["applicant_title"] = ""
	vars["technologies"] = strings.Join(p.Keywords, ", ")
	vars["skills"] = strings.Join(p.PreferredTech, ", ")
	vars["experience"] = ""
	vars["education"] = ""
	vars["achievements"] = ""
	vars["notable_projects"] = ""

	if p.CVData != nil {
		if p.CVData.Title != "" {
			vars["applicant_title"] = p.CVData.Title
		}
		if len(p.CVData.Technologies) > 0 {
			vars["technologies"] = strings.Join(p.CVData.Technologies, ", ")
		}
		if len(p.CVData.Skills) > 0 {
			vars["skills"] = strings.Join(p.CVData.Skills, ", ")
		}
		if len(p.CVData.Experience) > 0 {
			vars["experience"] = formatCoverLetterExperience(p.CVData.Experience)
			vars["notable_projects"] = extractNotableProjects(p.CVData.Experience)
		}
		if len(p.CVData.Education) > 0 {
			vars["education"] = formatCoverLetterEducation(p.CVData.Education)
		}
		if p.CVData.Summary != "" {
			vars["achievements"] = p.CVData.Summary
		}
	}

	prompt, err := llm.LoadPrompt("cover_letter", vars)
	if err != nil {
		return "", fmt.Errorf("load cover_letter prompt: %w", err)
	}

	body, err := client.Complete(
		"You are an expert career coach writing cover letters. Be professional, specific, and engaging.",
		prompt,
	)
	if err != nil {
		return "", fmt.Errorf("LLM generation: %w", err)
	}

	CacheCoverLetter(j.ID, body)

	return body, nil
}

func extractJobRequirements(description string) string {
	lines := strings.Split(description, "\n")
	inRequirements := false
	requirementLines := 0

	requirementPatterns := []string{
		"required:", "requirements:", "qualifications:",
		"must have", "must be", "required skills",
		"you have", "you bring", "we're looking for",
		"ideal candidate", "preferred qualifications",
	}

	var requirements []string

	for _, line := range lines {
		lowerLine := strings.ToLower(line)

		for _, pattern := range requirementPatterns {
			if strings.Contains(lowerLine, pattern) {
				inRequirements = true
				break
			}
		}

		if inRequirements && requirementLines < 15 {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				requirements = append(requirements, trimmed)
				requirementLines++
			}
		}

		if inRequirements && strings.Contains(lowerLine, "responsibilities") {
			break
		}
	}

	if len(requirements) == 0 {
		return truncate(description, 500)
	}

	return strings.Join(requirements, "\n")
}

func findMatchingTechnologies(jobDesc string, p *profile.Profile) string {
	desc := strings.ToLower(jobDesc)
	var matched []string

	allTech := append(p.PreferredTech, p.Keywords...)
	if p.CVData != nil {
		allTech = append(allTech, p.CVData.Technologies...)
	}

	seen := make(map[string]bool)
	for _, tech := range allTech {
		lowerTech := strings.ToLower(tech)
		if !seen[lowerTech] && strings.Contains(desc, lowerTech) {
			matched = append(matched, tech)
			seen[lowerTech] = true
		}
	}

	if len(matched) > 8 {
		matched = matched[:8]
	}

	return strings.Join(matched, ", ")
}

func extractNotableProjects(experiences []profile.Experience) string {
	var projects []string
	for i, exp := range experiences {
		if i >= 2 {
			break
		}
		if exp.Description != "" {
			desc := exp.Description
			if len(desc) > 150 {
				desc = desc[:150] + "..."
			}
			projects = append(projects, fmt.Sprintf("%s: %s", exp.Title, desc))
		}
	}
	return strings.Join(projects, "; ")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func formatCoverLetterExperience(experiences []profile.Experience) string {
	var parts []string
	for i, exp := range experiences {
		if i >= 3 {
			break
		}
		part := fmt.Sprintf("%s at %s", exp.Title, exp.Company)
		if exp.Duration != "" {
			part += fmt.Sprintf(" (%s)", exp.Duration)
		}
		if exp.Description != "" {
			desc := exp.Description
			if len(desc) > 200 {
				desc = desc[:200] + "..."
			}
			part += fmt.Sprintf(": %s", desc)
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, "; ")
}

func formatCoverLetterEducation(education []profile.Education) string {
	var parts []string
	for i, edu := range education {
		if i >= 2 {
			break
		}
		part := fmt.Sprintf("%s in %s from %s", edu.Degree, edu.Field, edu.Institution)
		if edu.Year != "" {
			part += fmt.Sprintf(" (%s)", edu.Year)
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, "; ")
}
