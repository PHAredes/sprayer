package cvwizard

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/theme"
)

func (m Model) updateSummary(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+enter":
			m.Step = StepExperience
			m.FocusIndex = 0
		case "esc":
			m.Step = StepInfo
			m.FocusIndex = 0
		case "backspace":
			if len(m.Summary) > 0 {
				m.Summary = m.Summary[:len(m.Summary)-1]
			}
		case "enter":
			m.Summary += "\n"
		default:
			if len(msg.String()) == 1 || msg.String() == " " {
				m.Summary += msg.String()
			}
		}
	}
	return m, nil
}

func (m Model) viewSummary() string {
	contentH := m.Height - 4
	if contentH < 1 {
		contentH = 1
	}

	bg := lipgloss.NewStyle().Background(theme.Background)

	hint := lipgloss.NewStyle().Foreground(theme.Muted).Render(
		"Used as context for email generation. 2-4 sentences is ideal.")

	bodyText := m.Summary + theme.CursorStyle.Render("|")

	textareaH := contentH - 4
	if textareaH < 3 {
		textareaH = 3
	}

	textarea := theme.CVTextareaStyle.
		Width(m.Width - 8).
		Height(textareaH).
		Padding(1, 2).
		Render(bodyText)

	body := lipgloss.JoinVertical(lipgloss.Left,
		hint,
		bg.Render(""),
		textarea,
	)

	return lipgloss.NewStyle().
		Background(theme.Background).
		Width(m.Width).
		Padding(1, 2).
		Render(body)
}
