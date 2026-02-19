package ui

import (
	"context"
	"fmt"
	"time"

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
	header  *Header
	footer  *Footer
	jobList *SimpleJobList
	scraper *ScraperView

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
	app.scraper = NewScraperView()

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
		// Add job to list incrementally
		a.jobs = append(a.jobs, msg.Job)
		a.filteredJobs = a.jobs
		a.jobList.SetJobs(a.jobs)
		return a, nil

	case AppScraperCompleteMsg:
		a.state = ModeJobs
		a.currentScraper = nil
		a.scraperCancel = nil
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
	}

	// Pass to job list component
	cmd := a.jobList.Update(msg)
	return a, cmd
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
		// Start monitoring in background
		go func() {
			for _ = range scraperInstance.Progress() {
				// Send progress message through program
				// This would need to be handled differently in bubbletea
			}
		}()

		// Collect results
		go func() {
			for _ = range scraperInstance.Results() {
				// Send job message through program
			}
		}()

		// Wait for completion
		go func() {
			// Wait for scraper completion
		}()

		return nil
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

	content := fmt.Sprintf("%s\n\n%s", title, GlobalKeys.FullHelp())
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
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
)

type keyMap struct {
	Quit key.Binding
	Help key.Binding
}

func (k keyMap) FullHelp() string {
	return fmt.Sprintf("Global Keys:\n  %s  %s\n  %s  %s",
		k.Quit.Help().Key, k.Quit.Help().Desc,
		k.Help.Help().Key, k.Help.Help().Desc,
	)
}

type jobKeyMap struct {
	Scrape  key.Binding
	Refresh key.Binding
	Up      key.Binding
	Down    key.Binding
	Select  key.Binding
}

func (k jobKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Scrape, k.Refresh, GlobalKeys.Help}
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
