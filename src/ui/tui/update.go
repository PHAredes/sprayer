package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"sprayer/src/ui/tui/compose"
	"sprayer/src/ui/tui/cvwizard"
	"sprayer/src/ui/tui/emails"
	"sprayer/src/ui/tui/filter"
	"sprayer/src/ui/tui/help"
	"sprayer/src/ui/tui/profiles"
	"sprayer/src/ui/tui/scraping"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window resize globally
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.height = ws.Height
		return m, nil
	}

	// Route to the active sub-screen
	switch m.viewState {
	case Filter:
		return m.updateFilter(msg)
	case Profiles:
		return m.updateProfiles(msg)
	case Help:
		return m.updateHelp(msg)
	case Scraping:
		return m.updateScraping(msg)
	case Emails:
		return m.updateEmails(msg)
	case Compose:
		return m.updateCompose(msg)
	case CVEntry, CVInfo, CVSummary, CVExperience, CVSkills, CVReview:
		return m.updateCVWizard(msg)
	default:
		return m.updateMain(msg)
	}
}

// ── Main screen (EmptyState / JobList) ─────────────────────────────────────

func (m Model) updateMain(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down", "↓":
			if len(m.jobs) > 0 {
				m.selectedIndex = min(m.selectedIndex+1, len(m.jobs)-1)
				m.viewState = JobList
			}
		case "k", "up", "↑":
			if len(m.jobs) > 0 {
				m.selectedIndex = max(m.selectedIndex-1, 0)
				m.viewState = JobList
			}
		case "s":
			m.prevState = m.viewState
			m.viewState = Scraping
			m.scrapingModel = scraping.New(9, m.width, m.height)
			return m, scraping.Tick()
		case "f":
			m.prevState = m.viewState
			m.viewState = Filter
			p := m.activeProfile
			m.filterModel = filter.New(
				strings.Join(p.Keywords, ", "),
				strings.Join(p.ExcludeKeywords, ", "),
				strings.Join(p.Locations, ", "),
				strings.Join(p.PreferredCompanies, ", "),
				fmt.Sprintf("%d", p.MinScore),
				m.width, m.height,
			)
		case "p":
			m.prevState = m.viewState
			m.viewState = Profiles
			m.profilesModel = profiles.New(m.allProfiles, 0, m.width, m.height)
		case "m":
			m.prevState = m.viewState
			m.viewState = Emails
			m.emailsModel = emails.New(m.emailDrafts, m.width, m.height)
		case "a":
			// apply placeholder
		case "?":
			m.prevState = m.viewState
			m.viewState = Help
			m.helpModel = help.New(m.width, m.height)
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

// ── Filter ─────────────────────────────────────────────────────────────────

func (m Model) updateFilter(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.filterModel, cmd = m.filterModel.Update(msg)

	// Check for transition messages
	if cmd != nil {
		resultMsg := cmd()
		switch resultMsg.(type) {
		case filter.FilterApplyMsg:
			applyMsg := resultMsg.(filter.FilterApplyMsg)
			m.activeProfile.Keywords = splitTrim(applyMsg.Keywords)
			m.activeProfile.ExcludeKeywords = splitTrim(applyMsg.Exclude)
			m.activeProfile.Locations = splitTrim(applyMsg.Locations)
			m.activeProfile.PreferredCompanies = splitTrim(applyMsg.Companies)
			m.viewState = m.prevState
			return m, nil
		case filter.FilterCancelMsg:
			m.viewState = m.prevState
			return m, nil
		}
	}

	return m, nil
}

// ── Profiles ───────────────────────────────────────────────────────────────

func (m Model) updateProfiles(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.profilesModel, cmd = m.profilesModel.Update(msg)

	if cmd != nil {
		resultMsg := cmd()
		switch resultMsg.(type) {
		case profiles.ProfileActivateMsg:
			actMsg := resultMsg.(profiles.ProfileActivateMsg)
			if actMsg.Index >= 0 && actMsg.Index < len(m.allProfiles) {
				m.activeProfile = m.allProfiles[actMsg.Index]
				m.profileName = m.activeProfile.Name
			}
			return m, nil
		case profiles.ProfileNewMsg:
			m.viewState = CVEntry
			m.cvWizardModel = cvwizard.New(m.width, m.height)
			return m, nil
		case profiles.ProfileBackMsg:
			m.viewState = m.prevState
			return m, nil
		}
	}

	return m, nil
}

// ── Help ───────────────────────────────────────────────────────────────────

func (m Model) updateHelp(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "esc", "?":
			m.viewState = m.prevState
			return m, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

// ── Scraping ───────────────────────────────────────────────────────────────

func (m Model) updateScraping(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.scrapingModel, cmd = m.scrapingModel.Update(msg)

	// Check for cancel
	if cmd != nil {
		resultMsg := cmd()
		switch resultMsg.(type) {
		case scraping.CancelMsg:
			m.viewState = m.prevState
			return m, nil
		case scraping.TickMsg:
			// Keep the tick running
			return m, scraping.Tick()
		}
	}

	// Handle tick messages - keep spinner going
	if _, ok := msg.(scraping.TickMsg); ok {
		return m, cmd
	}

	return m, cmd
}

// ── Emails ─────────────────────────────────────────────────────────────────

func (m Model) updateEmails(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.emailsModel, cmd = m.emailsModel.Update(msg)

	if cmd != nil {
		resultMsg := cmd()
		switch resultMsg.(type) {
		case emails.OpenComposeMsg:
			openMsg := resultMsg.(emails.OpenComposeMsg)
			if openMsg.Index >= 0 && openMsg.Index < len(m.emailDrafts) {
				e := m.emailDrafts[openMsg.Index]
				m.viewState = Compose
				m.composeModel = compose.New(
					e.From, e.To, e.Subject, e.Attach, e.Body,
					e.Company, 87, m.width, m.height,
				)
			}
			return m, nil
		case emails.BackMsg:
			m.viewState = m.prevState
			return m, nil
		case emails.DeleteMsg:
			delMsg := resultMsg.(emails.DeleteMsg)
			if delMsg.Index >= 0 && delMsg.Index < len(m.emailDrafts) {
				m.emailDrafts = append(m.emailDrafts[:delMsg.Index], m.emailDrafts[delMsg.Index+1:]...)
				m.emailsModel = emails.New(m.emailDrafts, m.width, m.height)
			}
			return m, nil
		}
	}

	return m, nil
}

// ── Compose ────────────────────────────────────────────────────────────────

func (m Model) updateCompose(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.composeModel, cmd = m.composeModel.Update(msg)

	if cmd != nil {
		resultMsg := cmd()
		switch resultMsg.(type) {
		case compose.BackMsg:
			m.viewState = Emails
			m.emailsModel = emails.New(m.emailDrafts, m.width, m.height)
			return m, nil
		case compose.SendMsg:
			// TODO: actual send
			m.viewState = Emails
			m.emailsModel = emails.New(m.emailDrafts, m.width, m.height)
			return m, nil
		}
	}

	return m, nil
}

// ── CV Wizard ──────────────────────────────────────────────────────────────

func (m Model) updateCVWizard(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.cvWizardModel, cmd = m.cvWizardModel.Update(msg)

	// Sync viewState with wizard step
	switch m.cvWizardModel.Step {
	case cvwizard.StepEntry:
		m.viewState = CVEntry
	case cvwizard.StepInfo:
		m.viewState = CVInfo
	case cvwizard.StepSummary:
		m.viewState = CVSummary
	case cvwizard.StepExperience:
		m.viewState = CVExperience
	case cvwizard.StepSkills:
		m.viewState = CVSkills
	case cvwizard.StepReview:
		m.viewState = CVReview
	}

	if cmd != nil {
		resultMsg := cmd()
		switch resultMsg.(type) {
		case cvwizard.DoneMsg:
			// Save profile and go back
			m.viewState = Profiles
			m.profilesModel = profiles.New(m.allProfiles, 0, m.width, m.height)
			return m, nil
		case cvwizard.CancelMsg:
			m.viewState = Profiles
			m.profilesModel = profiles.New(m.allProfiles, 0, m.width, m.height)
			return m, nil
		}
	}

	return m, nil
}

// ── Helpers ────────────────────────────────────────────────────────────────

func splitTrim(s string) []string {
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
