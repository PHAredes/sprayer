package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"sprayer/src/api/job"
	"sprayer/src/api/profile"
	"sprayer/src/ui/tui/compose"
	"sprayer/src/ui/tui/cvwizard"
	"sprayer/src/ui/tui/emails"
	"sprayer/src/ui/tui/filter"
	"sprayer/src/ui/tui/help"
	"sprayer/src/ui/tui/profiles"
	"sprayer/src/ui/tui/scraping"
)

type ViewState int

const (
	EmptyState ViewState = iota
	JobList
	Filter
	Profiles
	Help
	Scraping
	Emails
	Compose
	CVEntry
	CVInfo
	CVSummary
	CVExperience
	CVSkills
	CVReview
)

type Model struct {
	jobs          []job.Job
	selectedIndex int
	profileName   string
	viewState     ViewState
	prevState     ViewState // for back navigation
	width         int
	height        int

	// Sub-models
	filterModel   filter.Model
	profilesModel profiles.Model
	helpModel     help.Model
	scrapingModel scraping.Model
	emailsModel   emails.Model
	composeModel  compose.Model
	cvWizardModel cvwizard.Model

	// Data
	activeProfile profile.Profile
	allProfiles   []profile.Profile
	emailDrafts   []emails.Email
}

func NewModel() Model {
	dp := profile.NewDefaultProfile()
	return Model{
		jobs:          []job.Job{},
		selectedIndex: 0,
		profileName:   "Default",
		viewState:     EmptyState,
		prevState:     EmptyState,
		width:         80,
		height:        24,
		activeProfile: dp,
		allProfiles:   []profile.Profile{dp},
		emailDrafts:   sampleEmails(),
	}
}

func (m Model) SelectedIndex() int      { return m.selectedIndex }
func (m Model) ViewState() ViewState    { return m.viewState }
func (m Model) Jobs() []job.Job         { return m.jobs }
func (m *Model) SetJobs(jobs []job.Job) { m.jobs = jobs }

func (m Model) Init() tea.Cmd { return nil }

// sampleEmails provides demo email data.
func sampleEmails() []emails.Email {
	return []emails.Email{
		{
			Status: "draft", Company: "Vercel",
			Subject: "Application -- Sr. Software Engineer - Go",
			Date: "today", From: "you@example.com", To: "jobs@vercel.com",
			Attach: "resume.pdf",
			Body: "Hi Vercel team,\n\nI came across the Senior Software Engineer role and wanted to reach out directly -- the work you're doing on edge infrastructure and deployment pipelines is exactly the problem space I want to be in.\n\nI've spent the last 4 years writing Go for distributed systems, most recently shipping a multi-tenant job orchestration service with sub-100ms p99 latencies. Happy to share code samples or jump on a call.\n\n--\nyou@example.com - github.com/you",
		},
		{
			Status: "sent", Company: "PlanetScale",
			Subject: "Application -- Backend Engineer - Rust",
			Date: "yesterday", From: "you@example.com", To: "careers@planetscale.com",
			Attach: "resume.pdf", Body: "Dear PlanetScale team...",
		},
		{
			Status: "draft", Company: "Fly.io",
			Subject: "Application -- Platform Engineer - Remote",
			Date: "2d ago", From: "you@example.com", To: "jobs@fly.io",
			Attach: "resume.pdf", Body: "Hello Fly.io team...",
		},
		{
			Status: "sent", Company: "Stripe",
			Subject: "Application -- Go Infrastructure Engineer",
			Date: "3d ago", From: "you@example.com", To: "eng-recruiting@stripe.com",
			Attach: "resume.pdf", Body: "Hi Stripe team...",
		},
	}
}
