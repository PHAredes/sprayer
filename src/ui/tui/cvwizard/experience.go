package cvwizard

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/theme"
)

var expFields = []struct {
	Label       string
	Placeholder string
}{
	{"Role", "e.g. Senior Software Engineer"},
	{"Company", "e.g. Acme Corp"},
	{"Period", "e.g. 2021 - present"},
	{"Stack", "type and press enter to add tags"},
	{"Highlights", "type and press enter to add lines"},
}

func (m Model) updateExperience(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+enter":
			// Save current experience if it has content
			m.saveCurrentExp()
			m.Step = StepSkills
			m.FocusIndex = 0
		case "ctrl+n":
			// Save current and start new
			m.saveCurrentExp()
			m.CurrentExp = ExperienceEntry{}
			m.FocusIndex = 0
			m.CurrentTag = ""
		case "tab", "down":
			m.FocusIndex = (m.FocusIndex + 1) % len(expFields)
		case "shift+tab", "up":
			m.FocusIndex = (m.FocusIndex - 1 + len(expFields)) % len(expFields)
		case "esc":
			m.Step = StepSummary
			m.FocusIndex = 0
		case "enter":
			// For stack (3) and highlights (4), enter adds a tag/line
			if m.FocusIndex == 3 && m.CurrentTag != "" {
				m.CurrentExp.Stack = append(m.CurrentExp.Stack, m.CurrentTag)
				m.CurrentTag = ""
			} else if m.FocusIndex == 4 && m.CurrentTag != "" {
				m.CurrentExp.Highlights = append(m.CurrentExp.Highlights, m.CurrentTag)
				m.CurrentTag = ""
			}
		case "backspace":
			if m.FocusIndex == 3 || m.FocusIndex == 4 {
				if len(m.CurrentTag) > 0 {
					m.CurrentTag = m.CurrentTag[:len(m.CurrentTag)-1]
				} else if m.FocusIndex == 3 && len(m.CurrentExp.Stack) > 0 {
					m.CurrentExp.Stack = m.CurrentExp.Stack[:len(m.CurrentExp.Stack)-1]
				} else if m.FocusIndex == 4 && len(m.CurrentExp.Highlights) > 0 {
					m.CurrentExp.Highlights = m.CurrentExp.Highlights[:len(m.CurrentExp.Highlights)-1]
				}
			} else {
				field := m.expFieldPtr()
				if len(*field) > 0 {
					*field = (*field)[:len(*field)-1]
				}
			}
		default:
			if len(msg.String()) == 1 || msg.String() == " " {
				if m.FocusIndex == 3 || m.FocusIndex == 4 {
					m.CurrentTag += msg.String()
				} else {
					field := m.expFieldPtr()
					*field += msg.String()
				}
			}
		}
	}
	return m, nil
}

func (m *Model) saveCurrentExp() {
	if m.CurrentExp.Role != "" || m.CurrentExp.Company != "" {
		m.Experiences = append(m.Experiences, m.CurrentExp)
	}
}

func (m *Model) expFieldPtr() *string {
	switch m.FocusIndex {
	case 0:
		return &m.CurrentExp.Role
	case 1:
		return &m.CurrentExp.Company
	case 2:
		return &m.CurrentExp.Period
	}
	return &m.CurrentExp.Role
}

func (m Model) expFieldValue(idx int) string {
	switch idx {
	case 0:
		return m.CurrentExp.Role
	case 1:
		return m.CurrentExp.Company
	case 2:
		return m.CurrentExp.Period
	}
	return ""
}

func (m Model) viewExperience() string {
	contentH := m.Height - 4
	if contentH < 1 {
		contentH = 1
	}

	labelW := 12
	bg := lipgloss.NewStyle().Background(theme.Background)
	var rows []string

	for i, f := range expFields {
		focused := i == m.FocusIndex

		var lblStyle lipgloss.Style
		var rowBg lipgloss.Color
		var valStr string

		if focused {
			lblStyle = theme.CVFormLabelFocusedStyle
			rowBg = lipgloss.Color("#13102a")
		} else {
			lblStyle = theme.CVFormLabelStyle
			rowBg = theme.Background
		}

		if i < 3 {
			// Simple text fields
			val := m.expFieldValue(i)
			if focused {
				if val == "" {
					valStr = ""
				} else {
					valStr = lipgloss.NewStyle().Foreground(theme.Bright).Background(rowBg).Render(val)
				}
				valStr += theme.CursorStyle.Render("|")
			} else {
				if val == "" {
					valStr = theme.PlaceholderStyle.Background(rowBg).Render(f.Placeholder)
				} else {
					valStr = lipgloss.NewStyle().Foreground(theme.Text).Background(rowBg).Render(val)
				}
			}
		} else if i == 3 {
			// Stack tags
			valStr = m.renderTags(m.CurrentExp.Stack, focused, rowBg)
		} else if i == 4 {
			// Highlights
			valStr = m.renderHighlights(m.CurrentExp.Highlights, focused, rowBg)
		}

		label := lblStyle.Background(rowBg).Width(labelW).Align(lipgloss.Right).Render(f.Label)
		sep := lipgloss.NewStyle().Background(rowBg).Render("  ")

		var row string
		if focused {
			row = lipgloss.NewStyle().
				Background(rowBg).
				Width(m.Width - 6).
				PaddingLeft(1).
				PaddingRight(1).
				BorderStyle(lipgloss.NormalBorder()).
				BorderLeft(true).
				BorderForeground(theme.Purple).
				Render(label + sep + valStr)
		} else {
			row = lipgloss.NewStyle().
				Background(rowBg).
				Width(m.Width - 4).
				PaddingLeft(2).
				PaddingRight(1).
				Render(label + sep + valStr)
		}
		rows = append(rows, row)
	}

	// Added entries summary
	if len(m.Experiences) > 0 {
		rows = append(rows, bg.Render(""))
		rows = append(rows, bg.Foreground(theme.BorderColor).Width(m.Width-6).Render(strings.Repeat("-", m.Width-6)))

		addedLine := lipgloss.NewStyle().Foreground(theme.Muted).Render("Added: ")
		for _, exp := range m.Experiences {
			addedLine += lipgloss.NewStyle().Foreground(theme.Green).Render(
				fmt.Sprintf("  %s @ %s", exp.Role, exp.Company))
		}
		rows = append(rows, bg.PaddingLeft(2).Render(addedLine))
	}

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

func (m Model) renderTags(tags []string, focused bool, bg lipgloss.Color) string {
	var parts []string
	for _, t := range tags {
		parts = append(parts, theme.CVTagStyle.Render(" "+t+" "))
	}

	if focused {
		input := m.CurrentTag
		addTag := theme.CVTagAddFocusedStyle.Background(bg).Render("+ " + input)
		addTag += theme.CursorStyle.Render("|")
		parts = append(parts, addTag)
	} else {
		parts = append(parts, theme.CVTagAddStyle.Background(bg).Render("+ add"))
	}

	return strings.Join(parts, " ")
}

func (m Model) renderHighlights(highlights []string, focused bool, bg lipgloss.Color) string {
	var lines []string
	for _, h := range highlights {
		lines = append(lines, lipgloss.NewStyle().Foreground(theme.Text).Background(bg).Render("  · "+h))
	}
	if focused {
		input := m.CurrentTag
		addLine := theme.CVTagAddFocusedStyle.Background(bg).Render("+ " + input)
		addLine += theme.CursorStyle.Render("|")
		lines = append(lines, addLine)
	} else {
		lines = append(lines, theme.CVTagAddStyle.Background(bg).Render("  + add line"))
	}
	return strings.Join(lines, "\n")
}
