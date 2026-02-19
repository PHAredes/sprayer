package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/internal/job"
	"sprayer/internal/profile"
	"sprayer/internal/scraper"
)

// App represents the main TUI application following charm patterns
type App struct {
	// Core state
	state  AppMode
	width  int
	height int

	// Data
	jobs         []job.Job
	filteredJobs []job.Job
	profiles     []profile.Profile

	// Components
	header      *Header
	footer      *Footer
	jobList     *SimpleJobList
	jobDetail   *JobDetail
	filterView  *FilterView
	profileView *ProfileView
	scraper     *ScraperView
	helpView    help.Model

	// Dependencies
	store        *job.Store
	profileStore *profile.Store

	// Current scraping state
	currentScraper *scraper.IncrementalScraper
	scraperCancel  func()
}

// AppMode represents the current application mode
type AppMode int

const (
	ModeJobs AppMode = iota
	ModeJobDetail
	ModeFilters
	ModeProfiles
	ModeScraping
	ModeHelp
)

// NewApp creates a new TUI application
func NewApp() (*App, error) {
	// Initialize stores
	jobStore, err := job.NewStore()
	if err != nil {
		return nil, fmt.Errorf("failed to create job store: %w", err)
	}

	profileStore, err := profile.NewStore(jobStore.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create profile store: %w", err)
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

	app := &App{
		state:        ModeJobs,
		jobs:         jobs,
		filteredJobs: jobs,
		profiles:     profiles,
		store:        jobStore,
		profileStore: profileStore,
	}

	// Initialize components
	app.header = &Header{}
	app.footer = &Footer{}
	app.jobList = &SimpleJobList{jobs: jobs}
	app.jobDetail = NewJobDetail()
	app.filterView = NewFilterView(activeProfile)
	app.profileView = NewProfileView(profiles, activeProfile)
	app.scraper = NewScraperView()
	app.helpView = help.New()

	return app, nil
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, tea.EnterAltScreen)
	cmds = append(cmds, tea.SetWindowTitle("Sprayer - Job Application Tool"))

	// Initialize scraper view if in scraping mode
	if a.state == ModeScraping && a.scraper != nil {
		cmds = append(cmds, a.scraper.Init())
	}

	return tea.Batch(cmds...)
}

// Update handles messages and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.resizeComponents()
		return a, nil

	case tea.KeyMsg:
		// Global key handling
		switch {
		case key.Matches(msg, GlobalKeys.Quit):
			return a, tea.Quit
		case key.Matches(msg, GlobalKeys.Help):
			a.state = ModeHelp
			return a, nil
		}

		// State-specific key handling
		switch a.state {
		case ModeJobs:
			return a.updateJobs(msg)
		case ModeJobDetail:
			return a.updateJobDetail(msg)
		case ModeFilters:
			return a.updateFilters(msg)
		case ModeProfiles:
			return a.updateProfiles(msg)
		case ModeScraping:
			return a.updateScraping(msg)
		case ModeHelp:
			if msg.String() == "q" || msg.String() == "esc" || msg.String() == "?" {
				a.state = ModeJobs
			}
			return a, nil
		}

	case AppScraperProgressMsg:
		if a.scraper != nil {
			cmd := a.scraper.UpdateProgress(msg.Progress)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return a, tea.Batch(cmds...)

	case AppScraperJobMsg:
		// Add job to list incrementally and save to database
		a.jobs = append(a.jobs, msg.Job)
		a.filteredJobs = a.jobs
		a.jobList.SetJobs(a.jobs)

		// Save the individual job to database
		if err := a.store.Save([]job.Job{msg.Job}); err != nil {
			// Log error but don't stop the process
			fmt.Printf("Error saving job: %v\n", err)
		}

		return a, nil

	case AppScraperCompleteMsg:
		// Save all jobs and complete scraping
		a.state = ModeJobs
		a.currentScraper = nil
		a.scraperCancel = nil

		// Save all jobs to database and set last scrape time
		if err := a.store.Save(a.jobs); err != nil {
			fmt.Printf("Error saving jobs: %v\n", err)
		} else {
			fmt.Printf("Saved %d jobs from scraping\n", len(a.jobs))
		}

		return a, nil

	case AppScraperErrorMsg:
		// Handle scraper errors
		if a.scraper != nil {
			a.scraper.SetError(msg.Error)
		}
		return a, nil
	}

	// Update components
	if a.state == ModeJobs {
		cmd := a.jobList.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return a, tea.Batch(cmds...)
}

func (a *App) updateJobs(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, JobKeys.Scrape):
		return a.startScraping()
	case key.Matches(msg, JobKeys.Refresh):
		// Refresh job list
		jobs, _ := a.store.All()
		a.jobs = jobs
		a.filteredJobs = jobs
		a.jobList.SetJobs(jobs)
		return a, nil
	case key.Matches(msg, JobKeys.Filter):
		a.state = ModeFilters
		return a, nil
	case key.Matches(msg, JobKeys.Profiles):
		a.state = ModeProfiles
		return a, nil
	case key.Matches(msg, JobKeys.Select):
		// Select job and show detail
		selectedJob := a.jobList.SelectedJob()
		if selectedJob != nil && a.jobDetail != nil {
			a.jobDetail.SetJob(selectedJob)
			a.state = ModeJobDetail
		}
		return a, nil
	case key.Matches(msg, JobKeys.Apply):
		// Apply to selected job
		selectedJob := a.jobList.SelectedJob()
		if selectedJob != nil {
			fmt.Printf("Applying to job: %s @ %s\n", selectedJob.Title, selectedJob.Company)
			// TODO: Implement actual application logic
		}
		return a, nil
	}

	// Pass to job list component
	cmd := a.jobList.Update(msg)
	return a, cmd
}

func (a *App) updateJobDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		a.state = ModeJobs
		return a, nil
	case "a":
		// Apply to job
		if a.jobList.SelectedJob() != nil {
			fmt.Println("Applying to job...") // TODO: Implement application
		}
	}

	// Update job detail component
	if a.jobDetail != nil {
		cmd := a.jobDetail.Update(msg)
		return a, cmd
	}

	return a, nil
}

func (a *App) updateFilters(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		a.state = ModeJobs
		return a, nil
	case "enter":
		// Apply filters
		a.state = ModeJobs
		return a, nil
	}

	// Update filter view
	if a.filterView != nil {
		cmd := a.filterView.Update(msg)
		return a, cmd
	}

	return a, nil
}

func (a *App) updateProfiles(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		a.state = ModeJobs
		return a, nil
	case "enter":
		// Select profile
		a.state = ModeJobs
		return a, nil
	case "n":
		// New profile
		fmt.Println("Creating new profile...") // TODO: Implement profile creation
		return a, nil
	}

	// Update profile view
	if a.profileView != nil {
		cmd := a.profileView.Update(msg)
		return a, cmd
	}

	return a, nil
}

func (a *App) updateScraping(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		// Cancel scraping
		if a.scraperCancel != nil {
			a.scraperCancel()
		}
		a.state = ModeJobs
		return a, nil
	}

	// Update scraper view
	cmd := a.scraper.Update(msg)
	return a, cmd
}

func (a *App) startScraping() (tea.Model, tea.Cmd) {
	a.state = ModeScraping

	// Create context for scraping
	ctx, cancel := ContextWithScrapingTimeout()
	a.scraperCancel = cancel

	// Start incremental scraper
	scraper := scraper.NewIncrementalScraper(ctx, profile.NewDefaultProfile())
	a.currentScraper = scraper

	// Start scraping in background
	scraper.Start()

	// Monitor progress and results
	return a, tea.Batch(
		a.monitorScraper(scraper),
		a.scraper.Init(),
	)
}

func (a *App) monitorScraper(scraperInstance *scraper.IncrementalScraper) tea.Cmd {
	return func() tea.Msg {
		// Use a ticker-based approach to monitor scraper progress
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Check for progress updates
				select {
				case progress, ok := <-scraperInstance.Progress():
					if ok {
						return AppScraperProgressMsg{Progress: progress}
					}
				default:
				}

				// Check for job results
				select {
				case job, ok := <-scraperInstance.Results():
					if ok {
						return AppScraperJobMsg{Job: job}
					}
				default:
				}

				// Check if scraper is done
				select {
				case <-scraperInstance.Done():
					return AppScraperCompleteMsg{}
				default:
				}

			case <-scraperInstance.Done():
				return AppScraperCompleteMsg{}
			}
		}
	}
}

func (a *App) resizeComponents() {
	if a.header != nil {
		a.header.SetWidth(a.width)
	}
	if a.footer != nil {
		a.footer.SetWidth(a.width)
	}
	if a.jobList != nil {
		a.jobList.SetSize(a.width, a.height-2) // Account for header/footer
	}
	if a.scraper != nil {
		a.scraper.SetSize(a.width, a.height-2)
	}
}

// View renders the application
func (a *App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Loading..."
	}

	var content string
	switch a.state {
	case ModeJobs:
		content = a.jobList.View()
	case ModeJobDetail:
		if a.jobDetail != nil {
			content = a.jobDetail.View(a.width, a.height-2) // Account for header/footer
		} else {
			content = "No job selected"
		}
	case ModeFilters:
		if a.filterView != nil {
			content = a.filterView.View(a.width, a.height-2) // Account for header/footer
		} else {
			content = "Filter view not available"
		}
	case ModeProfiles:
		if a.profileView != nil {
			content = a.profileView.View(a.width, a.height-2) // Account for header/footer
		} else {
			content = "Profile view not available"
		}
	case ModeScraping:
		content = a.scraper.View()
	case ModeHelp:
		return a.renderHelp()
	}

	// Build full layout with header and footer
	return lipgloss.JoinVertical(
		lipgloss.Left,
		a.header.View(),
		content,
		a.footer.View(a.state),
	)
}

func (a *App) renderHelp() string {
	helpStyle := lipgloss.NewStyle().
		Background(Colors.Surface).
		Foreground(Colors.Text).
		Padding(2).
		Width(a.width).
		Height(a.height)

	title := lipgloss.NewStyle().
		Foreground(Colors.Primary).
		Bold(true).
		MarginBottom(1).
		Render("Sprayer - Help")

	// Use charm help component for better formatting
	jobHelp := a.helpView.View(JobKeys)
	globalHelp := a.helpView.View(GlobalKeys)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		"Job List Keys:",
		jobHelp,
		"",
		"Global Keys:",
		globalHelp,
	)

	return helpStyle.Render(content)
}

// Key maps following charm patterns
var (
	GlobalKeys = keyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}

	JobKeys = jobKeyMap{
		Scrape: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "scrape jobs"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Filter: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "filter"),
		),
		Profiles: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "profiles"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select job"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Apply: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "apply to job"),
		),
	}
)

type keyMap struct {
	Quit key.Binding
	Help key.Binding
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Help},
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

type jobKeyMap struct {
	Scrape   key.Binding
	Refresh  key.Binding
	Filter   key.Binding
	Profiles key.Binding
	Up       key.Binding
	Down     key.Binding
	Select   key.Binding
	Apply    key.Binding
}

func (k jobKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Scrape, k.Refresh, GlobalKeys.Help}
}

func (k jobKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{k.Scrape, k.Refresh},
	}
}

// Message types for scraper communication
type (
	AppScraperProgressMsg struct {
		Progress scraper.ScraperProgress
	}

	AppScraperJobMsg struct {
		Job job.Job
	}

	AppScraperCompleteMsg struct{}

	AppScraperErrorMsg struct {
		Error error
	}
)

// Helper functions
func ContextWithScrapingTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Minute)
}
