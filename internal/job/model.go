package job

import "time"

type Job struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Company      string    `json:"company"`
	Location     string    `json:"location"`
	Description  string    `json:"description"`
	URL          string    `json:"url"`
	Source       string    `json:"source"`
	PostedDate   time.Time `json:"posted_date"`
	Salary       string    `json:"salary,omitempty"`
	JobType      string    `json:"job_type,omitempty"`
	Email        string    `json:"email,omitempty"`
	Score        int       `json:"score"`
	HasTraps     bool      `json:"has_traps"`
	Traps        []string  `json:"traps,omitempty"`
	Applied      bool      `json:"applied"`
	AppliedDate  time.Time `json:"applied_date,omitempty"`
	ScratchEmail string    `json:"scratch_email,omitempty"`
	CustomCV     string    `json:"custom_cv,omitempty"`
}
