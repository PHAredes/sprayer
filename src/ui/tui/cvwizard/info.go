package cvwizard

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/theme"
)

var infoFields = []struct {
	Label       string
	Placeholder string
}{
	{"Name", "Your Name"},
	{"Job title", "e.g. Software Engineer"},
	{"Email", "you@example.com"},
	{"GitHub", "github.com/you"},
	{"LinkedIn", "linkedin.com/in/you"},
	{"Location", "City, Country"},
}

func (m Model) updateInfo(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.FocusIndex = (m.FocusIndex + 1) % len(infoFields)
		case "shift+tab", "up":
			m.FocusIndex = (m.FocusIndex - 1 + len(infoFields)) % len(infoFields)
		case "enter":
			m.Step = StepSummary
			m.FocusIndex = 0
		case "esc":
			m.Step = StepEntry
			m.FocusIndex = 0
		case "backspace":
			field := m.infoFieldPtr()
			if len(*field) > 0 {
				*field = (*field)[:len(*field)-1]
			}
		default:
			if len(msg.String()) == 1 || msg.String() == " " {
				field := m.infoFieldPtr()
				*field += msg.String()
			}
		}
	}
	return m, nil
}

func (m *Model) infoFieldPtr() *string {
	switch m.FocusIndex {
	case 0:
		return &m.Name
	case 1:
		return &m.JobTitle
	case 2:
		return &m.Email
	case 3:
		return &m.GitHub
	case 4:
		return &m.LinkedIn
	case 5:
		return &m.Location
	}
	return &m.Name
}

func (m Model) infoFieldValue(idx int) string {
	switch idx {
	case 0:
		return m.Name
	case 1:
		return m.JobTitle
	case 2:
		return m.Email
	case 3:
		return m.GitHub
	case 4:
		return m.LinkedIn
	case 5:
		return m.Location
	}
	return ""
}

func (m Model) viewInfo() string {
	contentH := m.Height - 4 // topbar, breadcrumb, statusbar
	if contentH < 1 {
		contentH = 1
	}

	labelW := 12
	var rows []string

	for i, f := range infoFields {
		focused := i == m.FocusIndex
		val := m.infoFieldValue(i)

		var lblStyle lipgloss.Style
		var rowBg lipgloss.Color
		var valStr string

		if focused {
			lblStyle = theme.CVFormLabelFocusedStyle
			rowBg = lipgloss.Color("#13102a")
			if val == "" {
				valStr = lipgloss.NewStyle().Foreground(theme.Bright).Background(rowBg).Render("")
			} else {
				valStr = lipgloss.NewStyle().Foreground(theme.Bright).Background(rowBg).Render(val)
			}
			valStr += theme.CursorStyle.Render("|")
		} else {
			lblStyle = theme.CVFormLabelStyle
			rowBg = theme.Background
			if val == "" {
				valStr = theme.PlaceholderStyle.Background(rowBg).Render(f.Placeholder)
			} else {
				valStr = lipgloss.NewStyle().Foreground(theme.Text).Background(rowBg).Render(val)
			}
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

	body := lipgloss.JoinVertical(lipgloss.Left, rows...)
	bodyH := lipgloss.Height(body)
	if bodyH < contentH {
		pad := lipgloss.NewStyle().Background(theme.Background).Width(m.Width).Height(contentH - bodyH).Render("")
		body = lipgloss.JoinVertical(lipgloss.Left, body, pad)
	}

	return lipgloss.NewStyle().
		Background(theme.Background).
		Width(m.Width).
		Padding(1, 2).
		Render(body)
}
