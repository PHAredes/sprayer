package cvwizard

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/theme"
)

var skillCategories = []string{"Languages", "Frameworks", "Databases", "Seniority"}

func (m Model) updateSkills(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+enter":
			m.Step = StepReview
			m.FocusIndex = 0
		case "tab", "down":
			m.FocusIndex = (m.FocusIndex + 1) % len(skillCategories)
			m.SkillInput = ""
		case "shift+tab", "up":
			m.FocusIndex = (m.FocusIndex - 1 + len(skillCategories)) % len(skillCategories)
			m.SkillInput = ""
		case "esc":
			m.Step = StepExperience
			m.FocusIndex = 0
		case "enter":
			if m.SkillInput != "" {
				tags := m.skillTagsPtr()
				*tags = append(*tags, strings.TrimSpace(m.SkillInput))
				m.SkillInput = ""
			}
		case "backspace":
			if len(m.SkillInput) > 0 {
				m.SkillInput = m.SkillInput[:len(m.SkillInput)-1]
			} else {
				tags := m.skillTagsPtr()
				if len(*tags) > 0 {
					*tags = (*tags)[:len(*tags)-1]
				}
			}
		default:
			if len(msg.String()) == 1 || msg.String() == " " {
				m.SkillInput += msg.String()
			}
		}
	}
	return m, nil
}

func (m *Model) skillTagsPtr() *[]string {
	switch m.FocusIndex {
	case 0:
		return &m.Languages
	case 1:
		return &m.Frameworks
	case 2:
		return &m.Databases
	case 3:
		return &m.Seniority
	}
	return &m.Languages
}

func (m Model) skillTags(idx int) []string {
	switch idx {
	case 0:
		return m.Languages
	case 1:
		return m.Frameworks
	case 2:
		return m.Databases
	case 3:
		return m.Seniority
	}
	return nil
}

func (m Model) viewSkills() string {
	contentH := m.Height - 4
	if contentH < 1 {
		contentH = 1
	}

	bg := lipgloss.NewStyle().Background(theme.Background)
	labelW := 12

	hint := lipgloss.NewStyle().Foreground(theme.Muted).Render(
		"These become your profile's keyword pool -- used for job matching and scoring.")
	hintBorder := bg.Foreground(theme.BorderColor).Width(m.Width - 8).Render(strings.Repeat("-", m.Width-8))

	var rows []string
	rows = append(rows, hint, hintBorder)

	for i, cat := range skillCategories {
		focused := i == m.FocusIndex
		tags := m.skillTags(i)

		var lblStyle lipgloss.Style
		var rowBg lipgloss.Color

		if focused {
			lblStyle = theme.CVFormLabelFocusedStyle
			rowBg = lipgloss.Color("#13102a")
		} else {
			lblStyle = theme.CVFormLabelStyle
			rowBg = theme.Background
		}

		label := lblStyle.Background(rowBg).Width(labelW).Align(lipgloss.Right).Render(cat)
		sep := lipgloss.NewStyle().Background(rowBg).Render("  ")

		// Tags
		var tagParts []string
		for _, t := range tags {
			tagParts = append(tagParts, theme.CVTagStyle.Render(" "+t+" "))
		}

		if focused {
			addTag := theme.CVTagAddFocusedStyle.Background(rowBg).Render("+ " + m.SkillInput)
			addTag += theme.CursorStyle.Render("|")
			tagParts = append(tagParts, addTag)
		} else {
			tagParts = append(tagParts, theme.CVTagAddStyle.Background(rowBg).Render("+ add"))
		}

		valStr := strings.Join(tagParts, " ")

		// Extra hint for seniority
		if i == 3 {
			valStr += lipgloss.NewStyle().Foreground(theme.Dim).Background(rowBg).Render(
				"  -- inferred from experience")
		}

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
		pad := bg.Width(m.Width).Height(contentH - bodyH).Render("")
		body = lipgloss.JoinVertical(lipgloss.Left, body, pad)
	}

	return lipgloss.NewStyle().
		Background(theme.Background).
		Width(m.Width).
		Padding(1, 2).
		Render(body)
}
