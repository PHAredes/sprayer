package cvwizard

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/theme"
)

func (m Model) updateReview(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg { return DoneMsg{ProfileName: m.ProfileName} }
		case "esc":
			m.Step = StepSkills
			m.FocusIndex = 0
		case "backspace":
			if len(m.ProfileName) > 0 {
				m.ProfileName = m.ProfileName[:len(m.ProfileName)-1]
			}
		default:
			if len(msg.String()) == 1 || msg.String() == " " {
				m.ProfileName += msg.String()
			}
		}
	}
	return m, nil
}

func (m Model) viewReview() string {
	contentH := m.Height - 4
	if contentH < 1 {
		contentH = 1
	}

	bg := lipgloss.NewStyle().Background(theme.Background)
	keyW := 14

	var rows []string

	// Name
	nameVal := lipgloss.NewStyle().Foreground(theme.Text).Render(m.Name)
	if m.JobTitle != "" || m.Location != "" {
		sub := " · "
		if m.JobTitle != "" {
			sub += m.JobTitle
		}
		if m.Location != "" {
			sub += " · " + m.Location
		}
		nameVal += lipgloss.NewStyle().Foreground(theme.Subtle).Render(sub)
	}
	rows = append(rows, m.reviewRow(keyW, "NAME", nameVal))

	// Summary
	summaryText := m.Summary
	if len(summaryText) > 60 {
		summaryText = summaryText[:57] + "..."
	}
	summaryVal := lipgloss.NewStyle().Foreground(theme.Subtle).Italic(true).Render(
		fmt.Sprintf("%q", summaryText))
	rows = append(rows, m.reviewRow(keyW, "SUMMARY", summaryVal))

	// Experience
	var expLines []string
	for _, exp := range m.Experiences {
		line := lipgloss.NewStyle().Foreground(theme.Bright).Render(exp.Role) +
			" @ " +
			lipgloss.NewStyle().Foreground(theme.Cyan).Render(exp.Company)
		if exp.Period != "" {
			line += "  " + lipgloss.NewStyle().Foreground(theme.Muted).Render(exp.Period)
		}
		expLines = append(expLines, line)
	}
	// Also show current entry if it has content
	if m.CurrentExp.Role != "" || m.CurrentExp.Company != "" {
		line := lipgloss.NewStyle().Foreground(theme.Bright).Render(m.CurrentExp.Role) +
			" @ " +
			lipgloss.NewStyle().Foreground(theme.Cyan).Render(m.CurrentExp.Company)
		if m.CurrentExp.Period != "" {
			line += "  " + lipgloss.NewStyle().Foreground(theme.Muted).Render(m.CurrentExp.Period)
		}
		expLines = append(expLines, line)
	}
	if len(expLines) == 0 {
		expLines = append(expLines, lipgloss.NewStyle().Foreground(theme.Muted).Render("none"))
	}
	rows = append(rows, m.reviewRow(keyW, "EXPERIENCE", strings.Join(expLines, "\n")))

	// Keywords
	kw := m.collectKeywords()
	var tagParts []string
	maxShow := 8
	for i, k := range kw {
		if i >= maxShow {
			tagParts = append(tagParts, lipgloss.NewStyle().Foreground(theme.Muted).Render(
				fmt.Sprintf("+%d more", len(kw)-maxShow)))
			break
		}
		tagParts = append(tagParts, theme.CVTagStyle.Render(" "+k+" "))
	}
	rows = append(rows, m.reviewRow(keyW, "KEYWORDS", strings.Join(tagParts, " ")))

	// Seniority
	var senParts []string
	for _, s := range m.Seniority {
		senParts = append(senParts, theme.CVTagStyle.Render(" "+s+" "))
	}
	if len(senParts) == 0 {
		senParts = append(senParts, lipgloss.NewStyle().Foreground(theme.Muted).Render("none"))
	}
	rows = append(rows, m.reviewRow(keyW, "SENIORITY", strings.Join(senParts, " ")))

	// Profile name (editable)
	nameInput := lipgloss.NewStyle().Foreground(theme.Bright).Render(m.ProfileName) +
		theme.CursorStyle.Render("|")
	rows = append(rows, m.reviewRow(keyW, "PROFILE NAME", nameInput))

	body := lipgloss.JoinVertical(lipgloss.Left, rows...)
	bodyH := lipgloss.Height(body)
	if bodyH < contentH {
		pad := bg.Width(m.Width).Height(contentH - bodyH).Render("")
		body = lipgloss.JoinVertical(lipgloss.Left, body, pad)
	}

	return lipgloss.NewStyle().
		Background(theme.Background).
		Width(m.Width).
		Padding(1, 2).
		Render(body)
}

func (m Model) reviewRow(keyW int, key, val string) string {
	bg := lipgloss.NewStyle().Background(theme.Background)

	k := theme.CVReviewKeyStyle.Width(keyW).Align(lipgloss.Right).Render(key)
	sep := bg.Render("  ")

	row := k + sep + val
	border := bg.Foreground(theme.BorderColor).Width(m.Width - 6).Render(strings.Repeat("-", m.Width-6))

	return lipgloss.JoinVertical(lipgloss.Left,
		bg.PaddingLeft(1).Width(m.Width-4).Render(row),
		border,
	)
}
