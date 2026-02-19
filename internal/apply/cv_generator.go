package apply

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"sprayer/internal/job"
	"sprayer/internal/llm"
	"sprayer/internal/profile"
)

type CVGenerator struct {
	client *llm.Client
	cache  map[string]*CachedCV
	mu     sync.RWMutex
}

type CachedCV struct {
	Content   string
	Generated time.Time
	JobID     string
}

func NewCVGenerator(client *llm.Client) *CVGenerator {
	return &CVGenerator{
		client: client,
		cache:  make(map[string]*CachedCV),
	}
}

func (g *CVGenerator) GenerateCustomCV(j *job.Job, p *profile.Profile) (string, error) {
	cacheKey := j.ID
	if g.client == nil {
		return "", fmt.Errorf("LLM client not available")
	}

	g.mu.RLock()
	if cached, ok := g.cache[cacheKey]; ok {
		if time.Since(cached.Generated) < 24*time.Hour {
			g.mu.RUnlock()
			return cached.Content, nil
		}
	}
	g.mu.RUnlock()

	cvData := p.CVData
	if cvData == nil && p.CVPath != "" {
		parser := profile.NewCVParser()
		var err error
		cvData, err = parser.ParseCVFromFile(p.CVPath)
		if err != nil {
			return "", fmt.Errorf("parse CV: %w", err)
		}
	}

	if cvData == nil {
		return "", fmt.Errorf("no CV data available for profile")
	}

	vars := map[string]string{
		"job_title":       j.Title,
		"company":         j.Company,
		"location":        j.Location,
		"job_description": truncate(j.Description, 3000),
		"applicant_name":  cvData.Name,
		"applicant_email": cvData.Email,
		"applicant_phone": cvData.Phone,
		"applicant_title": cvData.Title,
		"summary":         cvData.Summary,
		"technologies":    strings.Join(cvData.Technologies, ", "),
		"skills":          strings.Join(cvData.Skills, ", "),
		"experience":      formatExperience(cvData.Experience),
		"education":       formatEducation(cvData.Education),
	}

	prompt, err := llm.LoadPrompt("cv_custom", vars)
	if err != nil {
		return "", fmt.Errorf("load prompt: %w", err)
	}

	cvContent, err := g.client.Complete(
		"You are an expert CV/resume writer. Generate a tailored, professional CV that highlights relevant experience for the specific job. Be concise and impactful.",
		prompt,
	)
	if err != nil {
		return "", fmt.Errorf("LLM generation: %w", err)
	}

	g.mu.Lock()
	g.cache[cacheKey] = &CachedCV{
		Content:   cvContent,
		Generated: time.Now(),
		JobID:     j.ID,
	}
	g.mu.Unlock()

	return cvContent, nil
}

func (g *CVGenerator) GetCachedCV(jobID string) (string, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if cached, ok := g.cache[jobID]; ok {
		if time.Since(cached.Generated) < 24*time.Hour {
			return cached.Content, true
		}
	}
	return "", false
}

func (g *CVGenerator) ClearCache() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cache = make(map[string]*CachedCV)
}

func (g *CVGenerator) Available() bool {
	return g.client != nil && g.client.Available()
}

func SaveCustomCV(content, jobID, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	filename := fmt.Sprintf("custom_cv_%s_%d.txt", sanitize(jobID), time.Now().Unix())
	filepath := fmt.Sprintf("%s/%s", outputDir, filename)

	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write CV file: %w", err)
	}

	return filepath, nil
}

func formatExperience(experiences []profile.Experience) string {
	var parts []string
	for _, exp := range experiences {
		parts = append(parts, fmt.Sprintf("- %s at %s (%s): %s",
			exp.Title, exp.Company, exp.Duration, exp.Description))
	}
	return strings.Join(parts, "\n")
}

func formatEducation(education []profile.Education) string {
	var parts []string
	for _, edu := range education {
		parts = append(parts, fmt.Sprintf("- %s in %s from %s (%s)", edu.Degree, edu.Field, edu.Institution, edu.Year))
	}
	return strings.Join(parts, "\n")
}
