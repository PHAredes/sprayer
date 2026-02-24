package tui

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/joblist"
	"sprayer/src/ui/tui/theme"
)

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderTopBar(),
		m.renderContent(),
		m.renderStatusBar(),
	)
}

// ── Top bar — single row ──────────────────────────────────────────────────────

func (m Model) renderTopBar() string {
	on := func(fg lipgloss.Color) lipgloss.Style {
		return lipgloss.NewStyle().Background(theme.Surface).Foreground(fg)
	}

	left := on(theme.Subtle).Render("Profile: ") + on(theme.Cyan).Render(m.profileName)
	title := on(theme.Bright).Bold(true).Render("Sprayer")
	right := on(theme.Subtle).Render("Jobs: ") + on(theme.Yellow).Render(strconv.Itoa(len(m.jobs)))

	titleW := lipgloss.Width(title)
	sideW := (m.width - titleW) / 2
	rightW := m.width - sideW - titleW

	leftBlock := lipgloss.NewStyle().Background(theme.Surface).Width(sideW).PaddingLeft(2).Render(left)
	rightBlock := lipgloss.NewStyle().Background(theme.Surface).Width(rightW).PaddingRight(2).Align(lipgloss.Right).Render(right)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBlock, title, rightBlock)
}

// ── Content ───────────────────────────────────────────────────────────────────

func (m Model) renderContent() string {
	// Root level router for content. For now, we only have Job List (including Empty state).
	// Later, this handles help, filters, etc.
	switch m.viewState {
	case EmptyState, JobList:
		jm := joblist.Model{
			Jobs:          m.jobs,
			SelectedIndex: m.selectedIndex,
			Width:         m.width,
			Height:        m.height,
		}
		return jm.View()
	default:
		// Fallback for screens not yet implemented or managed at root.
		return lipgloss.NewStyle().
			Background(theme.Background).
			Width(m.width).
			Height(m.height - 2).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Screen [" + strconv.Itoa(int(m.viewState)) + "]")
	}
}

// ── Status bar — single row ───────────────────────────────────────────────────

func (m Model) renderStatusBar() string {
	keys := []string{"s", "f", "p", "m", "↑↓", "a", "?", "ctrl+c"}
	labels := []string{"scrape", "filter", "profiles", "emails", "navigate", "apply", "help", "quit"}

	// Footer kbd: same theme.Surface background as the bar — no tint.
	footerKbd := lipgloss.NewStyle().Background(theme.Surface).Foreground(theme.Cyan)
	sp := lipgloss.NewStyle().Background(theme.Surface).Foreground(theme.Subtle).Render(" ")

	line := ""
	for i, key := range keys {
		if i > 0 {
			line += theme.SepStyle.Render(" │ ")
		}
		line += footerKbd.Render(key) + sp + theme.StatusLabelStyle.Render(labels[i])
	}

	return lipgloss.NewStyle().Background(theme.Surface).Width(m.width).PaddingLeft(2).PaddingRight(2).Render(line)
}
