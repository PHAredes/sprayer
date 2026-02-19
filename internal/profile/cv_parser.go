package profile

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"sprayer/internal/job"
)

// CVData represents parsed CV information
type CVData struct {
	Name         string       `json:"name"`
	Email        string       `json:"email"`
	Phone        string       `json:"phone"`
	Location     string       `json:"location"`
	Title        string       `json:"title"`
	Summary      string       `json:"summary"`
	Technologies []string     `json:"technologies"`
	Experience   []Experience `json:"experience"`
	Education    []Education  `json:"education"`
	Skills       []string     `json:"skills"`
	Languages    []string     `json:"languages"`
}

// ScoreJob implements the CVMatcher interface
func (cv *CVData) ScoreJob(j *job.Job) int {
	score := 0
	jobText := strings.ToLower(j.Title + " " + j.Description)

	// Check technology matches
	for _, tech := range cv.Technologies {
		if strings.Contains(jobText, tech) {
			score += 10
		}
	}

	// Check skill matches
	for _, skill := range cv.Skills {
		if strings.Contains(jobText, skill) {
			score += 5
		}
	}

	// Check language matches
	for _, lang := range cv.Languages {
		if strings.Contains(jobText, lang) {
			score += 2
		}
	}

	// Check title match
	if cv.Title != "" && strings.Contains(jobText, strings.ToLower(cv.Title)) {
		score += 15
	}

	return score
}

// Experience represents work experience
type Experience struct {
	Company      string   `json:"company"`
	Title        string   `json:"title"`
	Duration     string   `json:"duration"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
}

// Education represents education background
type Education struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Field       string `json:"field"`
	Year        string `json:"year"`
}

// CVParser provides CV parsing functionality
type CVParser struct {
	techPatterns  []*regexp.Regexp
	skillPatterns []*regexp.Regexp
	langPatterns  []*regexp.Regexp
}

// NewCVParser creates a new CV parser
func NewCVParser() *CVParser {
	return &CVParser{
		techPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(go|golang)\b`),
			regexp.MustCompile(`(?i)\b(rust)\b`),
			regexp.MustCompile(`(?i)\b(python)\b`),
			regexp.MustCompile(`(?i)\b(javascript|js)\b`),
			regexp.MustCompile(`(?i)\b(typescript|ts)\b`),
			regexp.MustCompile(`(?i)\b(java)\b`),
			regexp.MustCompile(`(?i)\b(c\+\+)\b`),
			regexp.MustCompile(`(?i)\b(c#)\b`),
			regexp.MustCompile(`(?i)\b(react)\b`),
			regexp.MustCompile(`(?i)\b(vue)\b`),
			regexp.MustCompile(`(?i)\b(angular)\b`),
			regexp.MustCompile(`(?i)\b(node\.?js)\b`),
			regexp.MustCompile(`(?i)\b(django)\b`),
			regexp.MustCompile(`(?i)\b(flask)\b`),
			regexp.MustCompile(`(?i)\b(redis)\b`),
			regexp.MustCompile(`(?i)\b(postgres|postgresql)\b`),
			regexp.MustCompile(`(?i)\b(mongodb)\b`),
			regexp.MustCompile(`(?i)\b(docker)\b`),
			regexp.MustCompile(`(?i)\b(kubernetes|k8s)\b`),
			regexp.MustCompile(`(?i)\b(aws)\b`),
			regexp.MustCompile(`(?i)\b(gcp)\b`),
			regexp.MustCompile(`(?i)\b(azure)\b`),
		},
		skillPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(agile|scrum)\b`),
			regexp.MustCompile(`(?i)\b(git)\b`),
			regexp.MustCompile(`(?i)\b(ci/cd)\b`),
			regexp.MustCompile(`(?i)\b(devops)\b`),
			regexp.MustCompile(`(?i)\b(microservices)\b`),
			regexp.MustCompile(`(?i)\b(rest|restful)\b`),
			regexp.MustCompile(`(?i)\b(api)\b`),
			regexp.MustCompile(`(?i)\b(testing|unit tests|integration tests)\b`),
		},
		langPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(english)\b`),
			regexp.MustCompile(`(?i)\b(spanish)\b`),
			regexp.MustCompile(`(?i)\b(french)\b`),
			regexp.MustCompile(`(?i)\b(german)\b`),
			regexp.MustCompile(`(?i)\b(mandarin)\b`),
			regexp.MustCompile(`(?i)\b(portuguese)\b`),
		},
	}
}

// ParseCVFromText extracts CV data from text content
func (p *CVParser) ParseCVFromText(text string) (*CVData, error) {
	cv := &CVData{
		Technologies: []string{},
		Skills:       []string{},
		Languages:    []string{},
	}

	// Extract technologies
	for _, pattern := range p.techPatterns {
		if matches := pattern.FindAllString(text, -1); len(matches) > 0 {
			for _, match := range matches {
				tech := strings.ToLower(match)
				if !contains(cv.Technologies, tech) {
					cv.Technologies = append(cv.Technologies, tech)
				}
			}
		}
	}

	// Extract skills
	for _, pattern := range p.skillPatterns {
		if matches := pattern.FindAllString(text, -1); len(matches) > 0 {
			for _, match := range matches {
				skill := strings.ToLower(match)
				if !contains(cv.Skills, skill) {
					cv.Skills = append(cv.Skills, skill)
				}
			}
		}
	}

	// Extract languages
	for _, pattern := range p.langPatterns {
		if matches := pattern.FindAllString(text, -1); len(matches) > 0 {
			for _, match := range matches {
				lang := strings.ToLower(match)
				if !contains(cv.Languages, lang) {
					cv.Languages = append(cv.Languages, lang)
				}
			}
		}
	}

	return cv, nil
}

// ParseCVFromFile reads and parses a CV file
func (p *CVParser) ParseCVFromFile(filepath string) (*CVData, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CV file: %w", err)
	}

	return p.ParseCVFromText(string(content))
}

// GenerateProfileFromCV creates a profile from CV data
func GenerateProfileFromCV(cv *CVData, name string) Profile {
	prof := NewDefaultProfile()
	prof.Name = name
	prof.ID = fmt.Sprintf("cv_%d", time.Now().Unix())

	// Use technologies as keywords
	prof.Keywords = cv.Technologies

	// Use skills as preferred tech
	prof.PreferredTech = cv.Technologies

	// Set default preferences based on CV content
	if contains(cv.Technologies, "remote") || contains(cv.Skills, "remote") {
		prof.PreferRemote = true
	}

	// Add location if found
	if cv.Location != "" {
		prof.Locations = []string{cv.Location}
	}

	return prof
}

// CVBasedScorer scores jobs based on CV match
type CVBasedScorer struct {
	cv *CVData
}

// NewCVBasedScorer creates a new CV-based scorer
func NewCVBasedScorer(cv *CVData) *CVBasedScorer {
	return &CVBasedScorer{cv: cv}
}

// ScoreJob calculates a score based on CV match
func (s *CVBasedScorer) ScoreJob(j *job.Job) int {
	score := 0

	// Check technology matches
	jobText := strings.ToLower(j.Title + " " + j.Description)
	for _, tech := range s.cv.Technologies {
		if strings.Contains(jobText, tech) {
			score += 10
		}
	}

	// Check skill matches
	for _, skill := range s.cv.Skills {
		if strings.Contains(jobText, skill) {
			score += 5
		}
	}

	// Check language matches
	for _, lang := range s.cv.Languages {
		if strings.Contains(jobText, lang) {
			score += 2
		}
	}

	return score
}

// CVBasedFilter filters jobs based on CV match
type CVBasedFilter struct {
	cv       *CVData
	minScore int
}

// NewCVBasedFilter creates a new CV-based filter
func NewCVBasedFilter(cv *CVData, minScore int) *CVBasedFilter {
	return &CVBasedFilter{
		cv:       cv,
		minScore: minScore,
	}
}

// Filter applies CV-based filtering
func (f *CVBasedFilter) Filter(jobs []job.Job) []job.Job {
	var filtered []job.Job
	scorer := NewCVBasedScorer(f.cv)

	for _, j := range jobs {
		cvScore := scorer.ScoreJob(&j)
		if cvScore >= f.minScore {
			j.Score += cvScore // Add CV score to existing score
			filtered = append(filtered, j)
		}
	}

	return filtered
}

// Helper functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
