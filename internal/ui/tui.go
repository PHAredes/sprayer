package ui

import (
	"fmt"
	"strings"

	bubblekey "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"sprayer/internal/apply"
	"sprayer/internal/job"
	"sprayer/internal/llm"
	"sprayer/internal/profile"
	"sprayer/internal/scraper"
)

type state int

const (
	stateJobList state = iota
	stateJobDetail
	stateFilters
	stateProfiles
	stateScraping
	stateHelp
)

type Model struct {
	// Domain logic state
	jobs         []job.Job
	filteredJobs []job.Job
	profiles     []profile.Profile
	activeProfile profile.Profile
	filters      []job.Filter
	
	// UI state
	state        state
	width        int
	height       int
	cursor       int // index in filteredJobs
	scrollOffset int
	err          error
	statusMsg    string

	// Input components
	spinner      spinner.Model
	filterInput  textinput.Model
	viewport     viewport.Model

	// Services
	store        *job.Store
	profileStore *profile.Store
	llmClient    *llm.Client
}

func NewModel() (Model, error) {
	s, err := job.NewStore()
	if err != nil {
		return Model{}, err
	}

	pStore, err := profile.NewStore(s.DB)
	if err != nil {
		return Model{}, err
	}

	// Load initial data
	jobs, _ := s.All()
	profiles, _ := pStore.All()

	var active profile.Profile
	if len(profiles) > 0 {
		active = profiles[0]
	} else {
		// Create default profile if none exists
		active = profile.Profile{
			ID: "default", Name: "Default Profile", 
			Keywords: []string{"golang", "rust", "remote"},
		}
		pStore.Save(active)
		profiles = append(profiles, active)
	}

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(colorMauve)

	fi := textinput.New()
	fi.Placeholder = "Filter keywords..."
	fi.Cursor.Style = lipgloss.NewStyle().Foreground(colorMauve)

	return Model{
		store:        s,
		profileStore: pStore,
		jobs:         jobs,
		filteredJobs: jobs,
		profiles:     profiles,
		activeProfile: active,
		spinner:      sp,
		filterInput:  fi,
		llmClient:    llm.NewClient(),
		state:        stateJobList,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, textinput.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case bubblekey.Matches(msg, Keys.Quit):
			return m, tea.Quit
		case bubblekey.Matches(msg, Keys.Help):
			if m.state != stateHelp {
				m.state = stateHelp
			} else {
				m.state = stateJobList
			}
		case bubblekey.Matches(msg, Keys.Esc):
			if m.state != stateJobList {
				m.state = stateJobList
				m.err = nil
			}
		}

		// Mode-specific key handling
		switch m.state {
		case stateJobList:
			cmds = append(cmds, m.updateJobList(msg))
		case stateJobDetail:
			cmds = append(cmds, m.updateJobDetail(msg))
		case stateFilters:
			cmds = append(cmds, m.updateFilters(msg))
		case stateProfiles:
			cmds = append(cmds, m.updateProfiles(msg))
		case stateScraping:
			// logic handled in async message
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 10

	case jobsScrapedMsg:
		m.jobs = msg.jobs // merged result
		m.applyFilters()
		m.state = stateJobList
		m.statusMsg = fmt.Sprintf("Scraped %d jobs", len(msg.jobs))
		// Save to DB
		go m.store.Save(msg.jobs)

	case errMsg:
		m.err = msg.error
		m.state = stateJobList // return to list so user sees error in footer
	}

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateJobList(msg tea.KeyMsg) tea.Cmd {
	switch {
	case bubblekey.Matches(msg, Keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}
	case bubblekey.Matches(msg, Keys.Down):
		if m.cursor < len(m.filteredJobs)-1 {
			m.cursor++
		}
	case bubblekey.Matches(msg, Keys.Enter):
		if len(m.filteredJobs) > 0 {
			m.state = stateJobDetail
			m.viewport.SetContent(renderJobDetail(m.filteredJobs[m.cursor], m.activeProfile))
		}
	case bubblekey.Matches(msg, Keys.Scrape):
		m.state = stateScraping
		return m.startScraping()
	case bubblekey.Matches(msg, Keys.Filter):
		m.state = stateFilters
		m.filterInput.Focus()
	case bubblekey.Matches(msg, Keys.Profiles):
		m.state = stateProfiles
	case bubblekey.Matches(msg, Keys.Sort):
		// Toggle sort: Score -> Date
		// Naive toggle for now, could be cyclical
		m.filteredJobs = job.SortBy(job.ByDateDesc)(m.filteredJobs)
		m.statusMsg = "Sorted by Date"
	case bubblekey.Matches(msg, Keys.Apply):
		if len(m.filteredJobs) > 0 {
			j := m.filteredJobs[m.cursor]
			go func() {
				subject, body, err := apply.GenerateEmail(j, m.activeProfile, m.llmClient, "email_cold")
				if err != nil {
					// handle error
				} else {
					apply.Draft(j, m.activeProfile, subject, body)
				}
			}()
			m.statusMsg = fmt.Sprintf("Generating draft for %s...", j.Company)
		}
	case bubblekey.Matches(msg, Keys.Send):
		if len(m.filteredJobs) > 0 {
			j := m.filteredJobs[m.cursor]
			m.statusMsg = fmt.Sprintf("Generating and sending email to %s...", j.Company)
			go func() {
				// Generate content first
				subject, body, _ := apply.GenerateEmail(j, m.activeProfile, m.llmClient, "email_cold")
				// Then send directly
				err := apply.SendDirect(j.Email, subject, body, m.activeProfile.CVPath)
				if err != nil {
					// In a real app we'd dispatch an error Msg
				}
			}()
		}
	}
	return nil
}

func (m *Model) updateJobDetail(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return cmd
}

func (m *Model) updateFilters(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch {
	case bubblekey.Matches(msg, Keys.Enter):
		// Apply keywords filter
		kw := strings.Split(m.filterInput.Value(), ",")
		m.filters = []job.Filter{job.ByKeywords(kw)}
		m.applyFilters()
		m.state = stateJobList
		m.filterInput.Blur()
	case bubblekey.Matches(msg, Keys.Esc):
		m.filterInput.Blur()
		m.state = stateJobList
	}
	m.filterInput, cmd = m.filterInput.Update(msg)
	return cmd
}

func (m *Model) updateProfiles(msg tea.KeyMsg) tea.Cmd {
	// Simple cyclic switch for now
	switch {
	case bubblekey.Matches(msg, Keys.Down), bubblekey.Matches(msg, Keys.Up):
		// Find current index
		idx := 0
		for i, p := range m.profiles {
			if p.ID == m.activeProfile.ID {
				idx = i
				break
			}
		}
		idx = (idx + 1) % len(m.profiles)
		m.activeProfile = m.profiles[idx]
	case bubblekey.Matches(msg, Keys.Enter):
		m.state = stateJobList
	}
	return nil
}

func (m *Model) applyFilters() {
	// Start with all jobs -> dedup -> apply user filters -> sort by score
	pipeline := job.Pipe(
		job.Dedup(),
		job.Pipe(m.filters...),
		job.SortBy(job.ByScoreDesc),
	)
	m.filteredJobs = pipeline(m.jobs)
	m.cursor = 0
}

func (m *Model) startScraping() tea.Cmd {
	return func() tea.Msg {
		// Scrape All sources using profile keywords
		s := scraper.All(m.activeProfile.Keywords, "Remote") // default location for now
		jobs, err := s()
		if err != nil {
			return errMsg{err}
		}
		return jobsScrapedMsg{jobs}
	}
}

type jobsScrapedMsg struct {
	jobs []job.Job
}

type errMsg struct {
	error
}

func (m Model) View() string {
	if m.width == 0 {
		return "loading..."
	}

	switch m.state {
	case stateScraping:
		return fmt.Sprintf("\n\n   %s Scraping jobs... please wait.\n\n", m.spinner.View())
	case stateJobDetail:
		return m.viewJobDetail()
	case stateFilters:
		return m.viewFilters()
	case stateProfiles:
		return m.viewProfiles()
	case stateHelp:
		return m.viewHelp()
	default:
		return m.viewJobList()
	}
}
