package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"job-scraper/pkg/models"
	"job-scraper/pkg/scrapers"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Colors and styles
var (
	primaryColor    = lipgloss.Color("#FF6B9D")
	secondaryColor  = lipgloss.Color("#4ECDC4")
	accentColor     = lipgloss.Color("#45B7D1")
	warningColor    = lipgloss.Color("#FFA62B")
	successColor    = lipgloss.Color("#96CEB4")
	textColor       = lipgloss.Color("#F7F9F9")
	dimColor        = lipgloss.Color("#566573")

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(secondaryColor)

	statusStyle = lipgloss.NewStyle().
		Foreground(successColor).
		Italic(true)

	errorStyle = lipgloss.NewStyle().
		Foreground(warningColor).
		Bold(true)

	focusedStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)

	helpStyle = lipgloss.NewStyle().
		Foreground(dimColor)

	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

// Application states
type state int

const (
	stateJobList state = iota
	stateJobDetail
	stateFilter
	stateHelp
	stateScraping
	stateProfiles
	stateCVEditor
	stateSettings
)

// Key bindings
type keyMap struct {
	Select   key.Binding
	Back     key.Binding
	Filter   key.Binding
	Help     key.Binding
	Quit     key.Binding
	Export   key.Binding
	Scrape   key.Binding
	Profiles key.Binding
	CV       key.Binding
	Settings key.Binding
}

var keys = keyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "q"),
		key.WithHelp("esc/q", "back"),
	),
	Filter: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "filter jobs"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Export: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "export jobs"),
	),
	Scrape: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start scraping"),
	),
	Profiles: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "profiles"),
	),
	CV: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "cv editor"),
	),
	Settings: key.NewBinding(
		key.WithKeys(","),
		key.WithHelp(",", "settings"),
	),
}

// Implement help.KeyMap interface
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Back, k.Filter, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Back, k.Filter, k.Help},
		{k.Export, k.Scrape, k.Profiles, k.CV, k.Settings, k.Quit},
	}
}

// jobItem wraps models.Job to implement list.Item interface
type jobItem struct {
	models.Job
}

func (j jobItem) FilterValue() string {
	return j.Job.Title
}

func (j jobItem) Title() string {
	return fmt.Sprintf("%s - %s (Score: %d)", j.Job.Title, j.Job.Company, j.Job.Score)
}

func (j jobItem) Description() string {
	return fmt.Sprintf("%s | %s | %s", j.Job.Location, j.Job.Source, j.Job.PostedDate.Format("2006-01-02"))
}

// Main model
type model struct {
	state           state
	jobs            []models.Job
	filteredJobs    []models.Job
	selectedJob     *models.Job
	profiles        []models.Profile
	currentProfile  *models.Profile
	cvData          models.CVData
	list            list.Model
	jobViewport     viewport.Model
	filterInputs    []textinput.Model
	profileInputs   []textinput.Model
	cvInputs        []textinput.Model
	help            help.Model
	spinner         spinner.Model
	statusMessage   string
	err             error
	isScraping      bool
	db              *sql.DB
	configPath      string
}

// Messages for scraper integration
type scrapeStartedMsg struct{}
type scrapeCompleteMsg struct {
	jobs []models.Job
	err  error
}

// Start scraping command
func (m model) startScrapingCmd() tea.Cmd {
	return func() tea.Msg {
		scraper := scrapers.NewMockScraper("mock_scraper", m.db)
		jobs, err := scraper.Scrape()
		return scrapeCompleteMsg{jobs: jobs, err: err}
	}
}

// Initial model
func initialModel() model {
	db, err := models.InitDB()
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
	}

	// Load default profile
	profile := &models.Profile{
		ID:        "default",
		Name:      "Default Profile",
		Keywords:  []string{"rust", "golang", "compiler", "embedded"},
		Locations: []string{"remote"},
		MinScore:  80,
	}

	// Create default jobs list
	jobs := []models.Job{
		{
			ID:          "1",
			Title:       "Rust Developer",
			Company:     "Tech Corp",
			Location:    "Remote",
			Description: "Looking for experienced Rust developers...",
			URL:         "https://example.com/jobs/1",
			Source:      "Example",
			PostedDate:  time.Now(),
			Score:       85,
		},
	}

	// Create list items
	items := make([]list.Item, len(jobs))
	for i, job := range jobs {
		items[i] = jobItem{job}
	}

	// Configure list
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = titleStyle.Render("Job List")
	list.SetShowStatusBar(true)
	list.SetFilteringEnabled(true)

	// Configure job viewport
	jobViewport := viewport.New(80, 20)

	// Configure filter inputs
	filterInputs := make([]textinput.Model, 5)
	for i := range filterInputs {
		filterInputs[i] = textinput.New()
		filterInputs[i].CharLimit = 100
		filterInputs[i].Width = 40
	}
	filterInputs[0].Placeholder = "Keywords (comma separated)"
	filterInputs[1].Placeholder = "Minimum score (0-100)"
	filterInputs[2].Placeholder = "Location"
	filterInputs[3].Placeholder = "Company"
	filterInputs[4].Placeholder = "Date range (today/week/month/all)"

	// Configure spinner
	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		state:          stateJobList,
		jobs:           jobs,
		filteredJobs:   jobs,
		profiles:       []models.Profile{*profile},
		currentProfile: profile,
		list:           list,
		jobViewport:    jobViewport,
		filterInputs:   filterInputs,
		help:           help.New(),
		spinner:        s,
		statusMessage:  statusStyle.Render("Ready - Press 's' to start scraping"),
		db:             db,
		configPath:     filepath.Join(os.Getenv("HOME"), ".sprayer", "config.json"),
	}
}

// Init function
func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)

	case spinner.TickMsg:
		if m.isScraping {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case scrapeStartedMsg:
		m.statusMessage = statusStyle.Render("Starting job scraping...")
		return m, nil

	case scrapeCompleteMsg:
		m.isScraping = false
		if msg.err != nil {
			m.err = msg.err
			m.statusMessage = errorStyle.Render(fmt.Sprintf("Scraping failed: %v", msg.err))
		} else {
			m.jobs = msg.jobs
			m.filteredJobs = m.jobs
			m.updateJobList()
			m.state = stateJobList

			totalJobs := len(msg.jobs)
			m.statusMessage = statusStyle.Render(fmt.Sprintf("Scraping complete! Found %d jobs", totalJobs))

			// Check if we reached the target
			if totalJobs >= 300 {
				m.statusMessage += " 🎉 Target reached!"
			}
		}
		return m, nil
	}

	// Handle list updates
	if m.state == stateJobList {
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	// Handle viewport updates
	if m.state == stateJobDetail && m.selectedJob != nil {
		m.jobViewport, cmd = m.jobViewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

// Handle key messages
func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateJobList:
		return m.handleJobListKeyMsg(msg)
	case stateJobDetail:
		return m.handleJobDetailKeyMsg(msg)
	case stateFilter:
		return m.handleFilterKeyMsg(msg)
	case stateHelp:
		return m.handleHelpKeyMsg(msg)
	case stateScraping:
		return m.handleScrapingKeyMsg(msg)
	}
	return m, nil
}

// Handle job list key messages
func (m model) handleJobListKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Select):
		if selected := m.list.SelectedItem(); selected != nil {
			job := selected.(jobItem)
			m.selectedJob = &job.Job
			m.state = stateJobDetail
			m.jobViewport.SetContent(m.formatJobDetail(job.Job))
		}
		return m, nil

	case key.Matches(msg, keys.Filter):
		m.state = stateFilter
		// Focus the first filter input
		for i := range m.filterInputs {
			m.filterInputs[i].Blur()
		}
		m.filterInputs[0].Focus()
		return m, nil

	case key.Matches(msg, keys.Help):
		m.state = stateHelp
		return m, nil

	case key.Matches(msg, keys.Export):
		return m, m.exportJobs()

	case key.Matches(msg, keys.Scrape):
		m.state = stateScraping
		m.isScraping = true
		return m, m.startScrapingCmd()

	case key.Matches(msg, keys.Profiles):
		m.state = stateProfiles
		return m, nil

	case key.Matches(msg, keys.CV):
		m.state = stateCVEditor
		return m, nil

	case key.Matches(msg, keys.Settings):
		m.state = stateSettings
		return m, nil

	case key.Matches(msg, keys.Quit):
		return m, tea.Quit
	}
	return m, nil
}

// Handle job detail key messages
func (m model) handleJobDetailKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Back):
		m.state = stateJobList
		m.selectedJob = nil
		return m, nil
	}
	return m, nil
}

// Handle filter key messages
func (m model) handleFilterKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Find currently focused input
	focusedIndex := -1
	for i := range m.filterInputs {
		if m.filterInputs[i].Focused() {
			focusedIndex = i
			break
		}
	}

	switch {
	case key.Matches(msg, keys.Back):
		m.state = stateJobList
		// Blur all inputs when leaving filter mode
		for i := range m.filterInputs {
			m.filterInputs[i].Blur()
		}
		return m, nil

	case key.Matches(msg, keys.Select):
		// Enter key behavior
		if focusedIndex < len(m.filterInputs)-1 {
			// Move to next field
			m.filterInputs[focusedIndex].Blur()
			m.filterInputs[focusedIndex+1].Focus()
		} else {
			// Apply filters on last field
			m.applyFilter()
			m.state = stateJobList
			// Blur all inputs when leaving filter mode
			for i := range m.filterInputs {
				m.filterInputs[i].Blur()
			}
		}
		return m, nil

	case msg.Type == tea.KeyUp:
		// Move to previous field
		if focusedIndex > 0 {
			m.filterInputs[focusedIndex].Blur()
			m.filterInputs[focusedIndex-1].Focus()
		}
		return m, nil

	case msg.Type == tea.KeyDown || msg.Type == tea.KeyTab:
		// Move to next field
		if focusedIndex < len(m.filterInputs)-1 {
			m.filterInputs[focusedIndex].Blur()
			m.filterInputs[focusedIndex+1].Focus()
		}
		return m, nil

	case msg.Type == tea.KeyCtrlN:
		// Vim: next field
		if focusedIndex < len(m.filterInputs)-1 {
			m.filterInputs[focusedIndex].Blur()
			m.filterInputs[focusedIndex+1].Focus()
		}
		return m, nil

	case msg.Type == tea.KeyCtrlP:
		// Vim: previous field
		if focusedIndex > 0 {
			m.filterInputs[focusedIndex].Blur()
			m.filterInputs[focusedIndex-1].Focus()
		}
		return m, nil

	case msg.Type == tea.KeyHome:
		// Move to first field
		if focusedIndex >= 0 {
			m.filterInputs[focusedIndex].Blur()
		}
		m.filterInputs[0].Focus()
		return m, nil

	case msg.Type == tea.KeyEnd:
		// Move to last field
		if focusedIndex >= 0 {
			m.filterInputs[focusedIndex].Blur()
		}
		m.filterInputs[len(m.filterInputs)-1].Focus()
		return m, nil
	}

	// Update filter inputs
	var cmd tea.Cmd
	for i := range m.filterInputs {
		if m.filterInputs[i].Focused() {
			m.filterInputs[i], cmd = m.filterInputs[i].Update(msg)
		}
	}
	return m, cmd
}

// Handle help key messages
func (m model) handleHelpKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Back):
		m.state = stateJobList
		return m, nil
	}
	return m, nil
}

// Handle scraping key messages
func (m model) handleScrapingKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Back):
		m.state = stateJobList
		m.isScraping = false
		return m, nil
	}
	return m, nil
}

// Handle window size messages
func (m model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	h, v := docStyle.GetFrameSize()
	m.list.SetSize(msg.Width-h, msg.Height-v)
	m.jobViewport.Width = msg.Width - h
	m.jobViewport.Height = msg.Height - v - 3
	return m, nil
}

// Update job list
func (m *model) updateJobList() {
	items := make([]list.Item, len(m.filteredJobs))
	for i, job := range m.filteredJobs {
		items[i] = jobItem{job}
	}
	m.list.SetItems(items)
}

// Apply filter
func (m *model) applyFilter() {
	// Parse filter inputs
	keywords := strings.Split(m.filterInputs[0].Value(), ",")
	minScore := 0
	if m.filterInputs[1].Value() != "" {
		fmt.Sscanf(m.filterInputs[1].Value(), "%d", &minScore)
	}
	location := m.filterInputs[2].Value()
	company := m.filterInputs[3].Value()
	dateRange := m.filterInputs[4].Value()

	// Apply filters
	m.filteredJobs = m.jobs
	m.filteredJobs = m.filterByKeywords(m.filteredJobs, keywords)
	m.filteredJobs = m.filterByScore(m.filteredJobs, minScore)
	m.filteredJobs = m.filterByLocation(m.filteredJobs, location)
	m.filteredJobs = m.filterByCompany(m.filteredJobs, company)
	m.filteredJobs = m.filterByDateRange(m.filteredJobs, dateRange)

	m.updateJobList()
}

// Filter jobs by keywords
func (m model) filterByKeywords(jobs []models.Job, keywords []string) []models.Job {
	if len(keywords) == 0 {
		return jobs
	}

	var filtered []models.Job
	for _, job := range jobs {
		for _, keyword := range keywords {
			keyword = strings.TrimSpace(keyword)
			if strings.Contains(strings.ToLower(job.Title), strings.ToLower(keyword)) ||
				strings.Contains(strings.ToLower(job.Description), strings.ToLower(keyword)) {
				filtered = append(filtered, job)
				break
			}
		}
	}
	return filtered
}

// Filter jobs by minimum score
func (m model) filterByScore(jobs []models.Job, minScore int) []models.Job {
	if minScore <= 0 {
		return jobs
	}

	var filtered []models.Job
	for _, job := range jobs {
		if job.Score >= minScore {
			filtered = append(filtered, job)
		}
	}
	return filtered
}

// Filter jobs by location
func (m model) filterByLocation(jobs []models.Job, location string) []models.Job {
	if location == "" {
		return jobs
	}

	var filtered []models.Job
	for _, job := range jobs {
		if strings.Contains(strings.ToLower(job.Location), strings.ToLower(location)) {
			filtered = append(filtered, job)
		}
	}
	return filtered
}

// Filter jobs by company
func (m model) filterByCompany(jobs []models.Job, company string) []models.Job {
	if company == "" {
		return jobs
	}

	var filtered []models.Job
	for _, job := range jobs {
		if strings.Contains(strings.ToLower(job.Company), strings.ToLower(company)) {
			filtered = append(filtered, job)
		}
	}
	return filtered
}

// Filter jobs by date range
func (m model) filterByDateRange(jobs []models.Job, dateRange string) []models.Job {
	if dateRange == "" {
		return jobs
	}

	var filtered []models.Job
	now := time.Now()

	for _, job := range jobs {
		switch strings.ToLower(dateRange) {
		case "today":
			if job.PostedDate.Year() == now.Year() &&
				job.PostedDate.Month() == now.Month() &&
				job.PostedDate.Day() == now.Day() {
				filtered = append(filtered, job)
			}
		case "week":
			if job.PostedDate.After(now.AddDate(0, 0, -7)) {
				filtered = append(filtered, job)
			}
		case "month":
			if job.PostedDate.After(now.AddDate(0, -1, 0)) {
				filtered = append(filtered, job)
			}
		case "all":
			filtered = append(filtered, job)
		}
	}
	return filtered
}

// Export jobs to JSON
func (m model) exportJobs() tea.Cmd {
	return func() tea.Msg {
		filename := fmt.Sprintf("jobs_export_%s.json", time.Now().Format("20060102_150405"))
		filepath := filepath.Join("outputs", filename)

		// Create outputs directory if it doesn't exist
		os.MkdirAll("outputs", 0755)

		data, err := json.MarshalIndent(m.filteredJobs, "", "  ")
		if err != nil {
			log.Printf("Error marshaling jobs: %v", err)
			return nil
		}

		err = os.WriteFile(filepath, data, 0644)
		if err != nil {
			log.Printf("Error writing file: %v", err)
			return nil
		}

		log.Printf("Jobs exported to %s", filepath)
		return nil
	}
}

// Format job detail for viewport
func (m model) formatJobDetail(job models.Job) string {
	var sb strings.Builder
	sb.WriteString(headerStyle.Render(job.Title))
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("Company: %s\n", job.Company))
	sb.WriteString(fmt.Sprintf("Location: %s\n", job.Location))
	sb.WriteString(fmt.Sprintf("Source: %s\n", job.Source))
	sb.WriteString(fmt.Sprintf("Score: %d\n", job.Score))
	sb.WriteString(fmt.Sprintf("Posted: %s\n", job.PostedDate.Format("2006-01-02")))
	if job.Salary != "" {
		sb.WriteString(fmt.Sprintf("Salary: %s\n", job.Salary))
	}
	if job.JobType != "" {
		sb.WriteString(fmt.Sprintf("Job Type: %s\n", job.JobType))
	}
	sb.WriteString(fmt.Sprintf("URL: %s\n", job.URL))
	sb.WriteString("\n")
	sb.WriteString(headerStyle.Render("Description:"))
	sb.WriteString("\n")
	sb.WriteString(job.Description)
	return sb.String()
}

// View function
func (m model) View() string {
	switch m.state {
	case stateJobList:
		return m.renderJobListView()
	case stateJobDetail:
		return m.renderJobDetailView()
	case stateFilter:
		return m.renderFilterView()
	case stateHelp:
		return m.renderHelpView()
	case stateScraping:
		return m.renderScrapingView()
	case stateProfiles:
		return m.renderProfilesView()
	case stateCVEditor:
		return m.renderCVEditorView()
	case stateSettings:
		return m.renderSettingsView()
	default:
		return "Unknown state"
	}
}

// Render job list view
func (m model) renderJobListView() string {
	var sb strings.Builder
	sb.WriteString(m.list.View())
	sb.WriteString("\n")
	sb.WriteString(m.statusMessage)
	sb.WriteString("\n")
	// Custom footer instead of default help view
	sb.WriteString(helpStyle.Render("s:scrape • f:filter • e:export • p:profiles • c:cv • ?:help • ,:settings • ctrl+c:quit"))
	return docStyle.Render(sb.String())
}

// Render job detail view
func (m model) renderJobDetailView() string {
	var sb strings.Builder
	sb.WriteString(m.jobViewport.View())
	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("Press ESC to return to job list"))
	return docStyle.Render(sb.String())
}

// Render filter view
func (m model) renderFilterView() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Filter Jobs"))
	sb.WriteString("\n\n")

	// Field descriptions
	descriptions := []string{
		"Keywords (comma separated)",
		"Minimum score (0-100)",
		"Location",
		"Company", 
		"Date range (today/week/month/all)",
	}

	for i, input := range m.filterInputs {
		// Highlight focused field
		fieldLabel := fmt.Sprintf("%d. %s:", i+1, descriptions[i])
		if input.Focused() {
			fieldLabel = focusedStyle.Render(fieldLabel)
		} else {
			fieldLabel = headerStyle.Render(fieldLabel)
		}
		sb.WriteString(fieldLabel)
		sb.WriteString("\n  ")
		sb.WriteString(input.View())
		sb.WriteString("\n\n")
	}

	// Navigation help
	sb.WriteString(helpStyle.Render("↑/↓/ctrl+p/ctrl+n: Navigate • enter: Next field • esc: Cancel"))
	return docStyle.Render(sb.String())
}

// Render help view
func (m model) renderHelpView() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Sprayer Help"))
	sb.WriteString("\n\n")

	// Navigation help
	sb.WriteString(headerStyle.Render("Navigation:"))
	sb.WriteString("\n")
	sb.WriteString("  • enter - Select job / Next field (filter mode)")
	sb.WriteString("\n")
	sb.WriteString("  • esc/q - Go back / Cancel")
	sb.WriteString("\n")
	sb.WriteString("  • ctrl+c - Quit application")
	sb.WriteString("\n\n")

	// Main actions
	sb.WriteString(headerStyle.Render("Actions:"))
	sb.WriteString("\n")
	sb.WriteString("  • s - Start scraping for jobs")
	sb.WriteString("\n")
	sb.WriteString("  • f - Filter jobs")
	sb.WriteString("\n")
	sb.WriteString("  • e - Export jobs to JSON")
	sb.WriteString("\n")
	sb.WriteString("  • p - Manage profiles")
	sb.WriteString("\n")
	sb.WriteString("  • c - CV editor")
	sb.WriteString("\n")
	sb.WriteString("  • , - Settings")
	sb.WriteString("\n")
	sb.WriteString("  • ? - Show this help")
	sb.WriteString("\n\n")

	// Filter mode navigation
	sb.WriteString(headerStyle.Render("Filter Mode:"))
	sb.WriteString("\n")
	sb.WriteString("  • ↑/↓ or ctrl+p/ctrl/n - Navigate fields")
	sb.WriteString("\n")
	sb.WriteString("  • enter - Next field / Apply filters (last field)")
	sb.WriteString("\n")
	sb.WriteString("  • tab - Next field")
	sb.WriteString("\n")
	sb.WriteString("  • home - First field")
	sb.WriteString("\n")
	sb.WriteString("  • end - Last field")
	sb.WriteString("\n")
	sb.WriteString("  • esc - Cancel filtering")
	sb.WriteString("\n\n")

	// Filter fields description
	sb.WriteString(headerStyle.Render("Filter Fields:"))
	sb.WriteString("\n")
	sb.WriteString("  1. Keywords - Comma-separated terms (rust,compiler,embedded)")
	sb.WriteString("\n")
	sb.WriteString("  2. Min Score - 0-100 (80 for jobs scoring 80+)")
	sb.WriteString("\n")
	sb.WriteString("  3. Location - Remote, California, etc.")
	sb.WriteString("\n")
	sb.WriteString("  4. Company - Google, SpaceX, etc.")
	sb.WriteString("\n")
	sb.WriteString("  5. Date Range - today, week, month, all")
	sb.WriteString("\n\n")

	sb.WriteString(helpStyle.Render("Press ESC to return to job list"))
	return docStyle.Render(sb.String())
}

// Render scraping view
func (m model) renderScrapingView() string {
	var sb strings.Builder
	sb.WriteString(m.spinner.View())
	sb.WriteString(" Scraping jobs...\n")
	sb.WriteString(helpStyle.Render("Press ESC to cancel"))
	return docStyle.Render(sb.String())
}

// Render profiles view
func (m model) renderProfilesView() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Profiles"))
	sb.WriteString("\n\n")

	for i, profile := range m.profiles {
		if m.currentProfile != nil && profile.ID == m.currentProfile.ID {
			sb.WriteString(focusedStyle.Render(fmt.Sprintf("%d. %s (Active)", i+1, profile.Name)))
		} else {
			sb.WriteString(fmt.Sprintf("%d. %s", i+1, profile.Name))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("Press number to select profile • esc: back"))
	return docStyle.Render(sb.String())
}

// Render CV editor view
func (m model) renderCVEditorView() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render("CV Editor"))
	sb.WriteString("\n\n")
	sb.WriteString("CV editing feature coming soon...\n\n")
	sb.WriteString(helpStyle.Render("Press ESC to return to job list"))
	return docStyle.Render(sb.String())
}

// Render settings view
func (m model) renderSettingsView() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Settings"))
	sb.WriteString("\n\n")
	sb.WriteString("Settings management coming soon...\n\n")
	sb.WriteString(helpStyle.Render("Press ESC to return to job list"))
	return docStyle.Render(sb.String())
}

// Main function
func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}