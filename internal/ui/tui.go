package ui

import (
	"fmt"
	"strings"

	bubblekey "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"sprayer/internal/apply"
	"sprayer/internal/job"
	"sprayer/internal/llm"
	"sprayer/internal/parse"
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
	stateReviewEmail
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
	reviewInput  textarea.Model
	
	// Review state
	reviewJob    *job.Job
	reviewTraps  []string
	reviewSubject string

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

	ta := textarea.New()
	ta.Placeholder = "Email content..."
	ta.Focus()

	return Model{
		store:        s,
		profileStore: pStore,
		jobs:         jobs,
		filteredJobs: jobs,
		profiles:     profiles,
		activeProfile: active,
		spinner:      sp,
		filterInput:  fi,
		reviewInput:  ta,
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
		case stateReviewEmail:
			cmds = append(cmds, m.updateReviewEmail(msg))
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
		
	case emailGeneratedMsg:
		m.reviewSubject = msg.subject
		m.reviewInput.SetValue(msg.body)
		m.reviewInput.Focus()
		// If traps found, append to status?
		if len(m.reviewTraps) > 0 {
			m.statusMsg = fmt.Sprintf("WARNING: Traps detected: %v", m.reviewTraps)
		} else {
			m.statusMsg = "Draft generated. Edit and Ctrl+Enter to send."
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
			m.reviewJob = &j
			m.reviewTraps = parse.CheckForTraps(j.Description)
			
			m.state = stateReviewEmail
			m.statusMsg = fmt.Sprintf("Generating draft for %s...", j.Company)
			
			// Async generate and update textarea
			return func() tea.Msg {
				subject, body, err := apply.GenerateEmail(j, m.activeProfile, m.llmClient, "email_cold")
				if err != nil {
					return errMsg{err}
				}
				// Prepend subject to body for editing? Or handle separately.
				// Pop style: We just send body. Subject is usually separate.
				// For now, let's put Subject as first line?
				// Or better: Just body. Subject stays fixed or editable?
				// Let's make it Body only for simplicity, Subject is generated.
				// Or: "Subject: ... \n\n Body..." and parse it back?
				// Let's keep it simple: Edit Body.
				return emailGeneratedMsg{subject: subject, body: body}
			}
		}
	}
	return nil
}

type emailGeneratedMsg struct {
	subject, body string
} 

func (m *Model) updateReviewEmail(msg tea.KeyMsg) tea.Cmd {
	switch {
	case bubblekey.Matches(msg, Keys.Esc):
		m.state = stateJobList
		m.statusMsg = "Draft cancelled"
		return nil
	case msg.String() == "ctrl+enter":
		// Send
		if m.reviewJob != nil {
			go func(j job.Job, p profile.Profile, subject, body string) {
				apply.SendDirect(j.Email, subject, body, p.CVPath)
			}(*m.reviewJob, m.activeProfile, m.reviewSubject, m.reviewInput.Value())
			
			m.state = stateJobList
			m.statusMsg = fmt.Sprintf("Sent email to %s!", m.reviewJob.Company)
		}
		return nil
	}
	var cmd tea.Cmd
	m.reviewInput, cmd = m.reviewInput.Update(msg)
	return cmd
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
	case stateReviewEmail:
		return m.viewReviewEmail()
	case stateHelp:
		return m.viewHelp()
	default:
		return m.viewJobList()
	}
}
