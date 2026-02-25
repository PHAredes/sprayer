package cvwizard

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/theme"
)

func (m Model) updateEntry(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			m.ImportMode = false
			m.FocusIndex = 0
		case "right", "l":
			m.ImportMode = true
			m.FocusIndex = 0
		case "enter":
			m.Step = StepInfo
			m.FocusIndex = 0
		case "esc":
			return m, func() tea.Msg { return CancelMsg{} }
		case "backspace":
			if m.ImportMode && len(m.FilePath) > 0 {
				m.FilePath = m.FilePath[:len(m.FilePath)-1]
			}
		default:
			if m.ImportMode {
				if len(msg.String()) == 1 || msg.String() == " " {
					m.FilePath += msg.String()
				}
			}
		}
	}
	return m, nil
}

func (m Model) viewEntry() string {
	contentH := m.Height - 2
	if contentH < 1 {
		contentH = 1
	}

	bg := lipgloss.NewStyle().Background(theme.Background)

	// Description
	desc := lipgloss.NewStyle().Foreground(theme.Subtle).Render(
		"Sprayer uses your CV to match jobs by skills, seniority,\n" +
			"and experience -- and to personalize generated emails.")

	// Choice cards
	cardW := 30
	gap := 4

	manualIcon := theme.CVChoiceIconStyle.Render("*")
	manualLabel := lipgloss.NewStyle().Foreground(theme.Bright).Bold(true).Render("Fill manually")
	manualSub := lipgloss.NewStyle().Foreground(theme.Muted).Render("Step through a form.\nTakes ~3 minutes.")

	importIcon := theme.CVChoiceIconStyle.Render(">")
	importLabel := lipgloss.NewStyle().Foreground(theme.Bright).Bold(true).Render("Import file")
	importSub := lipgloss.NewStyle().Foreground(theme.Muted).Render(".pdf, .tex, .md, .txt\nExtracted via mods.")

	var manualStyle, importStyle lipgloss.Style
	if !m.ImportMode {
		manualStyle = theme.CVChoiceSelectedStyle.Width(cardW).Padding(1, 2).Align(lipgloss.Center)
		importStyle = theme.CVChoiceStyle.Width(cardW).Padding(1, 2).Align(lipgloss.Center)
		manualIcon = theme.CVChoiceIconSelectedStyle.Render("*")
	} else {
		manualStyle = theme.CVChoiceStyle.Width(cardW).Padding(1, 2).Align(lipgloss.Center)
		importStyle = theme.CVChoiceSelectedStyle.Width(cardW).Padding(1, 2).Align(lipgloss.Center)
		importIcon = theme.CVChoiceIconSelectedStyle.Render(">")
	}

	manualCard := manualStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, manualIcon, manualLabel, manualSub))
	importCard := importStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, importIcon, importLabel, importSub))

	cards := lipgloss.JoinHorizontal(lipgloss.Center,
		manualCard,
		lipgloss.NewStyle().Background(theme.Background).Width(gap).Render(""),
		importCard,
	)

	var parts []string
	parts = append(parts, desc, bg.Render(""), cards)

	// File path input (shown in import mode)
	if m.ImportMode {
		parts = append(parts, bg.Render(""))

		labelW := 12
		label := theme.CVFormLabelFocusedStyle.Width(labelW).Align(lipgloss.Right).Render("File path")
		val := lipgloss.NewStyle().Foreground(theme.Bright).Background(lipgloss.Color("#13102a")).Render(m.FilePath)
		cursor := theme.CursorStyle.Render("|")

		row := lipgloss.NewStyle().
			Background(lipgloss.Color("#13102a")).
			Width(cardW*2 + gap).
			PaddingLeft(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(theme.Purple).
			Render(label + "  " + val + cursor)

		parts = append(parts, row)

		hint := lipgloss.NewStyle().Foreground(theme.Muted).Render(
			"  mods will extract fields - you review every step before saving")
		parts = append(parts, hint)
	}

	block := lipgloss.JoinVertical(lipgloss.Center, parts...)

	return lipgloss.Place(
		m.Width, contentH,
		lipgloss.Center, lipgloss.Center,
		block,
		lipgloss.WithWhitespaceBackground(theme.Background),
	)
}
