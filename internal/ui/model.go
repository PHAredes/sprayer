package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"sprayer/internal/job"
	"sprayer/internal/profile"
	"sprayer/internal/scraper"
)

// Model is the main application model
type Model struct {
	// Core state
	state       AppState
	width       int
	height      int
	err         error
	statusMsg   string
	statusTimer int

	// Data
	jobs          []job.Job
	filteredJobs  []job.Job
	profiles      []profile.Profile
	activeProfile profile.Profile

	// Components
	jobList     *JobList
	jobDetail   *JobDetail
	profileView *ProfileView
	filterView  *FilterView
	statusBar   *StatusBar
	helpView    help.Model
	spinner     spinner.Model

	// Input components
	filterInput textinput.Model
	emailInput  textarea.Model

	// Dependencies
	store        *job.Store
	profileStore *profile.Store
	llmClient    interface{} // Will be properly typed

	// Navigation
	lastState AppState
	showHelp  bool
}

// AppState represents the current application state
type AppState int

const (
	StateList AppState = iota
	StateDetail
	StateFilters
	StateProfiles
	StateScraping
	StateHelp
	StateReview
	StateLoading
)

// NewModel creates a new TUI model
func NewModel() (Model, error) {
	// Initialize stores
	jobStore, err := job.NewStore()
	if err != nil {
		return Model{}, fmt.Errorf("failed to create job store: %w", err)
	}

	profileStore, err := profile.NewStore(jobStore.DB)
	if err != nil {
		return Model{}, fmt.Errorf("failed to create profile store: %w", err)
	}

	// Load data
	jobs, _ := jobStore.All()
	profiles, _ := profileStore.All()

	// Create default profile if none exists
	activeProfile := profile.NewDefaultProfile()
	if len(profiles) > 0 {
		activeProfile = profiles[0]
	} else {
		profileStore.Save(activeProfile)
		profiles = append(profiles, activeProfile)
	}

	// Initialize CHARM-style components
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(Colors.Accent)

	filterInput := textinput.New()
	filterInput.Placeholder = "Type to filter jobs..."
	filterInput.CharLimit = 100

	emailInput := textarea.New()
	emailInput.Placeholder = "Your application email..."
	emailInput.SetHeight(10)

	m := &Model{
		state:         StateList,
		jobs:          jobs,
		filteredJobs:  jobs,
		profiles:      profiles,
		activeProfile: activeProfile,
		store:         jobStore,
		profileStore:  profileStore,
		spinner:       sp,
		filterInput:   filterInput,
		emailInput:    emailInput,
		helpView:      help.New(),
	}

	// Initialize components
	m.jobList = NewJobList(m.filteredJobs)
	m.jobDetail = NewJobDetail()
	m.profileView = NewProfileView(profiles, activeProfile)
	m.filterView = NewFilterView(activeProfile)
	m.statusBar = NewStatusBar()

	// Apply initial filters
	m.applyProfileFilters()

	return *m, nil
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		textinput.Blink,
		textarea.Blink,
	)
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeComponents()

	case tea.KeyMsg:
		// Global key handling
		switch {
		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, Keys.Help):
			m.showHelp = !m.showHelp
			return m, nil
		case key.Matches(msg, Keys.Esc):
			if m.state != StateList {
				m.setState(StateList)
			}
			return m, nil
		}

		// state-specific key handling
		cmd := m.handleKeyMsg(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case jobsLoadedMsg:
		m.jobs = msg.jobs
		m.applyProfileFilters()
		m.setState(StateList)
		m.setStatus(fmt.Sprintf("Loaded %d jobs", len(msg.jobs)))

	case errorMsg:
		m.err = msg.err
		m.setStatus(fmt.Sprintf("Error: %v", msg.err))
		m.setState(StateList)
	}

	// Update current view component
	cmd := m.updateCurrentView(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		return m.renderLoading()
	}

	var content string
	switch m.state {
	case StateList:
		content = m.renderListView()
	case StateDetail:
		content = m.renderDetailView()
	case StateFilters:
		content = m.renderFilterView()
	case StateProfiles:
		content = m.renderProfileView()
	case StateScraping:
		content = m.renderScrapingView()
	case StateHelp:
		content = m.renderHelpView()
	case StateReview:
		content = m.renderReviewView()
	default:
		content = m.renderListView()
	}

	// Combine content with status bar
	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		m.statusBar.View(m.width, m.getStatusInfo()),
	)
}

// CHARM-style rendering methods
func (m *Model) renderLoading() string {
	return Styles.Loading.Render("Initializing Sprayer...")
}

func (m *Model) renderListView() string {
	header := m.renderHeader()
	list := m.jobList.View(m.width, m.height-3) // Account for header and status bar

	if m.showHelp {
		help := m.renderHelpOverlay()
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			help,
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(Colors.Surface),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, list)
}

func (m *Model) renderDetailView() string {
	return m.jobDetail.View(m.width, m.height-1)
}

func (m *Model) renderFilterView() string {
	return m.filterView.View(m.width, m.height-1)
}

func (m *Model) renderProfileView() string {
	return m.profileView.View(m.width, m.height-1)
}

func (m *Model) renderScrapingView() string {
	content := fmt.Sprintf("%s Scraping jobs with profile '%s'...",
		m.spinner.View(),
		m.activeProfile.Name,
	)

	return Styles.Scraping.
		Width(m.width).
		Height(m.height - 1).
		Render(content)
}

func (m *Model) renderHelpView() string {
	return m.helpView.View(Keys)
}

func (m *Model) renderReviewView() string {
	// Implementation for review view
	return "Review view - implementation pending"
}

func (m *Model) renderHeader() string {
	profileInfo := fmt.Sprintf("Profile: %s", m.activeProfile.Name)
	jobCount := fmt.Sprintf("Jobs: %d/%d", len(m.filteredJobs), len(m.jobs))

	left := Styles.HeaderText.Render(profileInfo)
	center := Styles.Title.Render("Sprayer")
	right := Styles.HeaderText.Render(jobCount)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		Styles.Header.Width(m.width/3).Render(left),
		Styles.Header.Width(m.width/3).Align(lipgloss.Center).Render(center),
		Styles.Header.Width(m.width/3).Align(lipgloss.Right).Render(right),
	)
}

func (m *Model) renderHelpOverlay() string {
	content := Styles.HelpBox.
		Width(m.width / 2).
		Height(m.height / 2).
		Render(m.helpView.View(Keys))

	return content
}

// Component update methods
func (m *Model) updateCurrentView(msg tea.Msg) tea.Cmd {
	switch m.state {
	case StateList:
		return m.jobList.Update(msg)
	case StateDetail:
		return m.jobDetail.Update(msg)
	case StateFilters:
		return m.filterView.Update(msg)
	case StateProfiles:
		return m.profileView.Update(msg)
	}
	return nil
}

// state management
func (m *Model) setState(state AppState) {
	m.lastState = m.state
	m.state = state
}

func (m *Model) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch m.state {
	case StateList:
		return m.handleListKeys(msg)
	case StateDetail:
		return m.handleDetailKeys(msg)
	case StateFilters:
		return m.handleFilterKeys(msg)
	case StateProfiles:
		return m.handleProfileKeys(msg)
	}
	return nil
}

// Key handlers for different States
func (m *Model) handleListKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Keys.Enter):
		if m.jobList.SelectedJob() != nil {
			m.jobDetail.SetJob(m.jobList.SelectedJob())
			m.setState(StateDetail)
		}
		return nil

	case key.Matches(msg, Keys.Scrape):
		m.setState(StateScraping)
		return m.startScraping()

	case key.Matches(msg, Keys.Filter):
		m.setState(StateFilters)
		return nil

	case key.Matches(msg, Keys.Profiles):
		m.setState(StateProfiles)
		return nil

	case key.Matches(msg, Keys.Sort):
		m.jobList.ToggleSort()
		return nil

	case key.Matches(msg, Keys.ClearFilter):
		m.clearFilters()
		return nil
	}
	return nil
}

func (m *Model) handleDetailKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Keys.Apply):
		// Start application process
		m.setState(StateReview)
		return nil
	case key.Matches(msg, Keys.Back), key.Matches(msg, Keys.Esc):
		m.setState(StateList)
		return nil
	}
	return nil
}

func (m *Model) handleFilterKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Keys.Enter):
		m.applyCustomFilters()
		m.setState(StateList)
		return nil
	case key.Matches(msg, Keys.Esc):
		m.setState(StateList)
		return nil
	}
	return nil
}

func (m *Model) handleProfileKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Keys.Enter):
		m.activeProfile = m.profileView.SelectedProfile()
		m.applyProfileFilters()
		m.setState(StateList)
		return nil
	case key.Matches(msg, Keys.Up), key.Matches(msg, Keys.Down):
		m.profileView.MoveCursor(msg.String() == "up")
		return nil
	case key.Matches(msg, Keys.Esc):
		m.setState(StateList)
		return nil
	}
	return nil
}

// Utility methods
func (m *Model) resizeComponents() {
	m.jobList.SetSize(m.width, m.height-3)
	m.jobDetail.SetSize(m.width, m.height-1)
	m.filterView.SetSize(m.width, m.height-1)
	m.profileView.SetSize(m.width, m.height-1)
}

func (m *Model) applyProfileFilters() {
	filters := m.activeProfile.GenerateFilters()
	m.filteredJobs = job.Pipe(filters...)(m.jobs)
	m.jobList.SetJobs(m.filteredJobs)
	m.setStatus(fmt.Sprintf("Applied %d filters from profile '%s'", len(filters), m.activeProfile.Name))
}

func (m *Model) applyCustomFilters() {
	// Apply custom filters from filter view
	// Implementation depends on filter view design
	m.applyProfileFilters()
}

func (m *Model) clearFilters() {
	m.filteredJobs = m.jobs
	m.jobList.SetJobs(m.filteredJobs)
	m.setStatus("Cleared all filters")
}

func (m *Model) startScraping() tea.Cmd {
	return func() tea.Msg {
		// Use new profile-based scraper
		s := scraper.ProfileBasedScraper(m.activeProfile)
		jobs, err := s()
		if err != nil {
			return errorMsg{err}
		}
		return jobsLoadedMsg{jobs}
	}
}

func (m *Model) setStatus(msg string) {
	m.statusMsg = msg
	m.statusTimer = 60 // Show for ~3 seconds at 20fps
}

func (m *Model) getStatusInfo() StatusBarInfo {
	return StatusBarInfo{
		Message:     m.statusMsg,
		Profile:     m.activeProfile.Name,
		JobCount:    len(m.filteredJobs),
		TotalJobs:   len(m.jobs),
		StatusTimer: m.statusTimer,
	}
}

// Message types
type jobsLoadedMsg struct {
	jobs []job.Job
}

type errorMsg struct {
	err error
}

// StatusInfo contains information for the status bar
type StatusInfo struct {
	Message     string
	Profile     string
	JobCount    int
	TotalJobs   int
	statusTimer int
}
