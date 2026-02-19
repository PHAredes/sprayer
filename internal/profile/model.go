package profile

import (
	"fmt"
	"strings"
	"time"

	"sprayer/internal/job"
)

// Profile represents a person-specific application profile.
// Links to a CV variant, cover letter template, and keywords for matching.
type Profile struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Keywords     []string `json:"keywords"`
	CVPath       string   `json:"cv_path"`
	CoverPath    string   `json:"cover_path"`
	ContactEmail string   `json:"contact_email"`
	PreferRemote bool     `json:"prefer_remote"`
	Locations    []string `json:"locations"`

	// Dynamic filtering configuration
	MinScore        int         `json:"min_score"`
	MaxScore        int         `json:"max_score"`
	ExcludeTraps    bool        `json:"exclude_traps"`
	MustHaveEmail   bool        `json:"must_have_email"`
	JobTypes        []string    `json:"job_types"`        // "full-time", "contract", "part-time", "internship"
	SeniorityLevels []string    `json:"seniority_levels"` // "junior", "mid", "senior", "staff", "principal"
	SalaryRange     SalaryRange `json:"salary_range"`
	ExcludeKeywords []string    `json:"exclude_keywords"`

	// Technology preferences
	PreferredTech []string `json:"preferred_tech"`
	AvoidTech     []string `json:"avoid_tech"`

	// Company preferences
	PreferredCompanies []string `json:"preferred_companies"`
	AvoidCompanies     []string `json:"avoid_companies"`

	// Date filtering
	PostedAfter  *time.Time `json:"posted_after"`
	PostedBefore *time.Time `json:"posted_before"`

	// Scoring weights (0-100)
	ScoringWeights ScoringWeights `json:"scoring_weights"`

	// CV-based data
	CVData     *CVData `json:"cv_data,omitempty"`
	CVMinScore int     `json:"cv_min_score,omitempty"` // Minimum CV match score
}

type SalaryRange struct {
	Min      int    `json:"min"`
	Max      int    `json:"max"`
	Currency string `json:"currency"`
}

type ScoringWeights struct {
	TechMatch      int `json:"tech_match"`
	SeniorityMatch int `json:"seniority_match"`
	LocationMatch  int `json:"location_match"`
	CompanyMatch   int `json:"company_match"`
	SalaryMatch    int `json:"salary_match"`
	RemoteMatch    int `json:"remote_match"`
}

// NewDefaultProfile creates a profile with sensible defaults
func NewDefaultProfile() Profile {
	now := time.Now()
	return Profile{
		ID:              "default",
		Name:            "Default",
		Keywords:        []string{"golang", "rust", "remote"},
		MinScore:        0,
		MaxScore:        100,
		ExcludeTraps:    true,
		MustHaveEmail:   false,
		JobTypes:        []string{"full-time", "contract"},
		SeniorityLevels: []string{"mid", "senior", "staff"},
		PostedAfter:     &now, // Default to jobs posted today
		ScoringWeights:  DefaultScoringWeights(),
	}
}

func DefaultScoringWeights() ScoringWeights {
	return ScoringWeights{
		TechMatch:      30,
		SeniorityMatch: 20,
		LocationMatch:  15,
		CompanyMatch:   10,
		SalaryMatch:    15,
		RemoteMatch:    10,
	}
}

// GenerateFilters creates job filters based on profile preferences
func (p *Profile) GenerateFilters() []job.Filter {
	var filters []job.Filter

	// Keyword filters
	if len(p.Keywords) > 0 {
		filters = append(filters, job.ByKeywords(p.Keywords))
	}

	// Exclude keywords
	if len(p.ExcludeKeywords) > 0 {
		filters = append(filters, job.ExcludeKeywords(p.ExcludeKeywords))
	}

	// Location filters
	if len(p.Locations) > 0 {
		filters = append(filters, job.ByLocations(p.Locations))
	}

	// Company filters
	if len(p.PreferredCompanies) > 0 {
		filters = append(filters, job.ByCompanies(p.PreferredCompanies))
	}

	if len(p.AvoidCompanies) > 0 {
		filters = append(filters, job.ExcludeCompanies(p.AvoidCompanies))
	}

	// Score range
	if p.MinScore > 0 || p.MaxScore < 100 {
		filters = append(filters, job.ByScoreRange(p.MinScore, p.MaxScore))
	}

	// Email requirement
	if p.MustHaveEmail {
		filters = append(filters, job.HasEmail())
	}

	// Trap exclusion
	if p.ExcludeTraps {
		filters = append(filters, job.ExcludeTraps())
	}

	// Remote preference
	if p.PreferRemote {
		filters = append(filters, job.RemotePreferred())
	}

	// Seniority level
	if len(p.SeniorityLevels) > 0 {
		filters = append(filters, job.BySeniorityLevels(p.SeniorityLevels))
	}

	// Technology preferences
	if len(p.PreferredTech) > 0 {
		filters = append(filters, job.ByTechnologies(p.PreferredTech))
	}

	if len(p.AvoidTech) > 0 {
		filters = append(filters, job.ExcludeTechnologies(p.AvoidTech))
	}

	// Date filtering
	if p.PostedAfter != nil {
		filters = append(filters, job.PostedAfter(*p.PostedAfter))
	}

	if p.PostedBefore != nil {
		filters = append(filters, job.PostedBefore(*p.PostedBefore))
	}

	// CV-based filtering
	if p.CVData != nil && p.CVMinScore > 0 {
		filters = append(filters, job.ByCVMatch(p.CVData, p.CVMinScore))
	}

	return filters
}

// CalculateJobScore calculates a custom score for a job based on profile preferences
func (p *Profile) CalculateJobScore(j *job.Job) int {
	score := 0
	maxScore := 0

	// Technology matching
	if len(p.PreferredTech) > 0 {
		maxScore += p.ScoringWeights.TechMatch
		titleDesc := strings.ToLower(j.Title + " " + j.Description)
		for _, tech := range p.PreferredTech {
			if strings.Contains(titleDesc, strings.ToLower(tech)) {
				score += p.ScoringWeights.TechMatch / len(p.PreferredTech)
				break
			}
		}
	}

	// Seniority matching
	if len(p.SeniorityLevels) > 0 {
		maxScore += p.ScoringWeights.SeniorityMatch
		titleLower := strings.ToLower(j.Title)
		for _, level := range p.SeniorityLevels {
			if strings.Contains(titleLower, level) {
				score += p.ScoringWeights.SeniorityMatch
				break
			}
		}
	}

	// Location matching
	if len(p.Locations) > 0 && p.PreferRemote {
		maxScore += p.ScoringWeights.LocationMatch
		if strings.Contains(strings.ToLower(j.Location), "remote") {
			score += p.ScoringWeights.LocationMatch
		}
	}

	// Company matching
	if len(p.PreferredCompanies) > 0 {
		maxScore += p.ScoringWeights.CompanyMatch
		for _, company := range p.PreferredCompanies {
			if strings.Contains(strings.ToLower(j.Company), strings.ToLower(company)) {
				score += p.ScoringWeights.CompanyMatch
				break
			}
		}
	}

	// Remote matching
	if p.PreferRemote {
		maxScore += p.ScoringWeights.RemoteMatch
		if strings.Contains(strings.ToLower(j.Location), "remote") {
			score += p.ScoringWeights.RemoteMatch
		}
	}

	// Normalize to 0-100 scale
	if maxScore > 0 {
		return (score * 100) / maxScore
	}

	return 50 // Default neutral score
}

// GetFilterSummary returns a human-readable summary of active filters
func (p *Profile) GetFilterSummary() string {
	var parts []string

	if len(p.Keywords) > 0 {
		parts = append(parts, fmt.Sprintf("keywords: %s", strings.Join(p.Keywords, ", ")))
	}

	if len(p.Locations) > 0 {
		parts = append(parts, fmt.Sprintf("locations: %s", strings.Join(p.Locations, ", ")))
	}

	if p.MinScore > 0 {
		parts = append(parts, fmt.Sprintf("min score: %d", p.MinScore))
	}

	if p.ExcludeTraps {
		parts = append(parts, "no traps")
	}

	if p.MustHaveEmail {
		parts = append(parts, "has email")
	}

	if p.PreferRemote {
		parts = append(parts, "remote preferred")
	}

	if len(p.SeniorityLevels) > 0 {
		parts = append(parts, fmt.Sprintf("levels: %s", strings.Join(p.SeniorityLevels, ", ")))
	}

	if len(parts) == 0 {
		return "no filters"
	}

	return strings.Join(parts, " • ")
}
