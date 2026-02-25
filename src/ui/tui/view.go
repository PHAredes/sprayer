package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/joblist"
	"sprayer/src/ui/tui/statusbar"
	"sprayer/src/ui/tui/theme"
)

func (m Model) View() string {
	parts := []string{m.renderTopBar()}

	// Some screens have a modal top bar (breadcrumbs)
	if modal := m.renderModalTopBar(); modal != "" {
		parts = append(parts, modal)
	}

	parts = append(parts, m.renderContent())
	parts = append(parts, m.renderStatusBar())

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// ── Top bar ────────────────────────────────────────────────────────────────

func (m Model) renderTopBar() string {
	on := func(fg lipgloss.Color) lipgloss.Style {
		return lipgloss.NewStyle().Background(theme.Surface).Foreground(fg)
	}

	var left, title, right string

	switch m.viewState {
	case Profiles:
		left = on(theme.Subtle).Render("Profiles")
		title = on(theme.Bright).Bold(true).Render("Sprayer")
		right = on(theme.Subtle).Render("up/dn nav  |  enter select  |  n new  |  esc back")

	case Help:
		left = on(theme.Subtle).Render("")
		title = on(theme.Bright).Bold(true).Render("Sprayer -- Help")
		right = on(theme.Subtle).Render("")

	case Emails:
		left = on(theme.Subtle).Render("Emails")
		title = on(theme.Bright).Bold(true).Render("Sprayer")
		right = on(theme.Subtle).Render("Drafts: ") +
			on(theme.Yellow).Render(strconv.Itoa(m.emailsModel.DraftCount())) +
			on(theme.Subtle).Render("  ·  Sent: ") +
			on(theme.Yellow).Render(strconv.Itoa(m.emailsModel.SentCount()))

	case Compose:
		left = on(theme.Subtle).Render("Compose")
		title = on(theme.Bright).Bold(true).Render("Sprayer")
		right = on(theme.Subtle).Render(m.composeModel.Company) +
			on(theme.Subtle).Render("  ·  Score: ") +
			on(theme.Yellow).Render(strconv.Itoa(m.composeModel.Score))

	case CVEntry:
		left = on(theme.Subtle).Render("Profiles")
		title = on(theme.Bright).Bold(true).Render("Sprayer -- New Profile")
		right = on(theme.Subtle).Render("")

	default:
		left = on(theme.Subtle).Render("Profile: ") + on(theme.Cyan).Render(m.profileName)
		title = on(theme.Bright).Bold(true).Render("Sprayer")
		right = on(theme.Subtle).Render("Jobs: ") + on(theme.Yellow).Render(strconv.Itoa(len(m.jobs)))
	}

	titleW := lipgloss.Width(title)
	sideW := (m.width - titleW) / 2
	rightW := m.width - sideW - titleW

	leftBlock := lipgloss.NewStyle().Background(theme.Surface).Width(sideW).PaddingLeft(2).Render(left)
	rightBlock := lipgloss.NewStyle().Background(theme.Surface).Width(rightW).PaddingRight(2).Align(lipgloss.Right).Render(right)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBlock, title, rightBlock)
}

// ── Modal top bar (for filter and CV wizard) ──────────────────────────────

func (m Model) renderModalTopBar() string {
	switch m.viewState {
	case Filter:
		m.filterModel.Width = m.width
		return m.filterModel.ModalTopBar()
	case CVInfo, CVSummary, CVExperience, CVSkills, CVReview:
		m.cvWizardModel.Width = m.width
		return m.cvWizardModel.ModalTopBar()
	}
	return ""
}

// ── Content ───────────────────────────────────────────────────────────────

func (m Model) renderContent() string {
	switch m.viewState {
	case EmptyState, JobList:
		jm := joblist.Model{
			Jobs:          m.jobs,
			SelectedIndex: m.selectedIndex,
			Width:         m.width,
			Height:        m.height,
		}
		return jm.View()

	case Filter:
		m.filterModel.Width = m.width
		m.filterModel.Height = m.height
		return m.filterModel.View()

	case Profiles:
		m.profilesModel.Width = m.width
		m.profilesModel.Height = m.height
		return m.profilesModel.View()

	case Help:
		m.helpModel.Width = m.width
		m.helpModel.Height = m.height
		return m.helpModel.View()

	case Scraping:
		m.scrapingModel.Width = m.width
		m.scrapingModel.Height = m.height
		return m.scrapingModel.View()

	case Emails:
		m.emailsModel.Width = m.width
		m.emailsModel.Height = m.height
		return m.emailsModel.View()

	case Compose:
		m.composeModel.Width = m.width
		m.composeModel.Height = m.height
		return m.composeModel.View()

	case CVEntry, CVInfo, CVSummary, CVExperience, CVSkills, CVReview:
		m.cvWizardModel.Width = m.width
		m.cvWizardModel.Height = m.height
		return m.cvWizardModel.View()

	default:
		return lipgloss.NewStyle().
			Background(theme.Background).
			Width(m.width).
			Height(m.height - 2).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Screen [" + fmt.Sprintf("%d", int(m.viewState)) + "]")
	}
}

// ── Status bar ────────────────────────────────────────────────────────────

func (m Model) renderStatusBar() string {
	switch m.viewState {
	case Filter:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "tab", Label: "next"},
			{Key: "enter", Label: "apply"},
			{Key: "esc", Label: "cancel"},
			{Key: "?", Label: "help"},
			{Key: "ctrl+c", Label: "quit"},
		}, "Profile: "+m.profileName+"  |  No traps")

	case Profiles:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "up/dn", Label: "select"},
			{Key: "enter", Label: "choose"},
			{Key: "n", Label: "new"},
			{Key: "esc", Label: "back"},
			{Key: "ctrl+c", Label: "quit"},
		}, "")

	case Help:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "esc", Label: "close"},
			{Key: "ctrl+c", Label: "quit"},
		}, "")

	case Scraping:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "esc", Label: "cancel"},
			{Key: "?", Label: "help"},
			{Key: "ctrl+c", Label: "quit"},
		}, "")

	case Emails:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "up/dn", Label: "navigate"},
			{Key: "enter", Label: "open"},
			{Key: "g", Label: "generate new"},
			{Key: "d", Label: "delete"},
			{Key: "esc", Label: "back"},
		}, "mods + pop")

	case Compose:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "ctrl+s", Label: "send via pop"},
			{Key: "ctrl+r", Label: "regenerate"},
			{Key: "tab", Label: "next field"},
			{Key: "e", Label: "edit body"},
			{Key: "esc", Label: "back"},
		}, fmt.Sprintf("%d words", m.composeModel.WordCount()))

	case CVEntry:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "<->", Label: "choose"},
			{Key: "enter", Label: "continue"},
			{Key: "esc", Label: "cancel"},
		}, "")

	case CVInfo:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "tab/dn", Label: "next field"},
			{Key: "up", Label: "prev field"},
			{Key: "enter", Label: "continue"},
			{Key: "esc", Label: "back"},
		}, m.cvWizardModel.StepDots())

	case CVSummary:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "ctrl+enter", Label: "continue"},
			{Key: "esc", Label: "back"},
		}, m.cvWizardModel.StepDots())

	case CVExperience:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "tab", Label: "next field"},
			{Key: "ctrl+n", Label: "add another"},
			{Key: "ctrl+enter", Label: "next section"},
			{Key: "esc", Label: "back"},
		}, m.cvWizardModel.StepDots())

	case CVSkills:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "tab", Label: "next field"},
			{Key: "enter", Label: "confirm tag"},
			{Key: "backspace", Label: "remove"},
			{Key: "ctrl+enter", Label: "continue"},
			{Key: "esc", Label: "back"},
		}, m.cvWizardModel.StepDots())

	case CVReview:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "enter", Label: "save profile"},
			{Key: "tab", Label: "edit section"},
			{Key: "esc", Label: "back"},
		}, m.cvWizardModel.StepDots())

	default:
		return statusbar.Render(m.width, []statusbar.Binding{
			{Key: "s", Label: "scrape"},
			{Key: "f", Label: "filter"},
			{Key: "p", Label: "profiles"},
			{Key: "m", Label: "emails"},
			{Key: "up/dn", Label: "navigate"},
			{Key: "a", Label: "apply"},
			{Key: "?", Label: "help"},
			{Key: "ctrl+c", Label: "quit"},
		}, "")
	}
}
