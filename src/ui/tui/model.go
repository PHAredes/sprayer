package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"sprayer/src/api/job"
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
	width         int
	height        int
}

func NewModel() Model {
	return Model{
		jobs:          []job.Job{},
		selectedIndex: 0,
		profileName:   "Default",
		viewState:     EmptyState,
		width:         80,
		height:        24,
	}
}

func (m *Model) SelectedIndex() int     { return m.selectedIndex }
func (m *Model) ViewState() ViewState   { return m.viewState }
func (m *Model) Jobs() []job.Job        { return m.jobs }
func (m *Model) SetJobs(jobs []job.Job) { m.jobs = jobs }

func (m Model) Init() tea.Cmd { return nil }