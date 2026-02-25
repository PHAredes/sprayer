package profiles

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/src/api/profile"
	"sprayer/src/ui/tui/theme"
)

type Model struct {
	Profiles      []profile.Profile
	SelectedIndex int
	ActiveIndex   int
	Width         int
	Height        int
}

func New(profiles []profile.Profile, active int, w, h int) Model {
	if len(profiles) == 0 {
		profiles = []profile.Profile{profile.NewDefaultProfile()}
	}
	return Model{
		Profiles:      profiles,
		SelectedIndex: 0,
		ActiveIndex:   active,
		Width:         w,
		Height:        h,
	}
}

// ProfileActivateMsg signals profile activation.
type ProfileActivateMsg struct{ Index int }

// ProfileNewMsg signals new profile creation.
type ProfileNewMsg struct{}

// ProfileBackMsg signals going back.
type ProfileBackMsg struct{}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.SelectedIndex < len(m.Profiles)-1 {
				m.SelectedIndex++
			}
		case "k", "up":
			if m.SelectedIndex > 0 {
				m.SelectedIndex--
			}
		case "enter":
			m.ActiveIndex = m.SelectedIndex
			return m, func() tea.Msg { return ProfileActivateMsg{Index: m.SelectedIndex} }
		case "n":
			return m, func() tea.Msg { return ProfileNewMsg{} }
		case "esc":
			return m, func() tea.Msg { return ProfileBackMsg{} }
		}
	}
	return m, nil
}

func (m Model) View() string {
	contentH := m.Height - 2
	if contentH < 1 {
		contentH = 1
	}

	listW := m.Width * 24 / 100
	if listW < 20 {
		listW = 20
	}
	detailW := m.Width - listW

	left := m.renderList(listW, contentH)
	right := m.renderDetail(detailW, contentH)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m Model) renderList(w, h int) string {
	bg := lipgloss.NewStyle().Background(theme.Background)

	header := bg.
		Foreground(theme.Muted).
		Width(w).
		PaddingLeft(2).
		PaddingRight(1).
		Bold(true).
		Render("PROFILES")

	headerBorder := bg.Foreground(theme.BorderColor).Width(w).Render(strings.Repeat("-", w))

	var items []string
	items = append(items, header, headerBorder)

	for i, p := range m.Profiles {
		selected := i == m.SelectedIndex
		isActive := i == m.ActiveIndex

		namePrefix := "  "
		if selected {
			namePrefix = "> "
		}

		badge := ""
		if isActive {
			badge = " " + theme.ActiveBadgeStyle.Render("active")
		}

		kwSummary := summarizeKeywords(p.Keywords, w-8)

		var nameStyle, subStyle lipgloss.Style
		if selected {
			nameStyle = lipgloss.NewStyle().Foreground(theme.Purple).Bold(true).Background(lipgloss.Color("#13102a"))
			subStyle = lipgloss.NewStyle().Foreground(theme.Muted).Background(lipgloss.Color("#13102a"))
		} else {
			nameStyle = lipgloss.NewStyle().Foreground(theme.Bright).Background(theme.Background)
			subStyle = lipgloss.NewStyle().Foreground(theme.Muted).Background(theme.Background)
		}

		var rowBg lipgloss.Color
		if selected {
			rowBg = lipgloss.Color("#13102a")
		} else {
			rowBg = theme.Background
		}

		name := nameStyle.Render(namePrefix+p.Name) + badge
		sub := subStyle.Render("  " + kwSummary)

		rowContent := lipgloss.JoinVertical(lipgloss.Left, name, sub)

		item := lipgloss.NewStyle().
			Background(rowBg).
			Width(w).
			PaddingLeft(1).
			Render(rowContent)

		if selected {
			item = lipgloss.NewStyle().
				Background(rowBg).
				Width(w).
				BorderStyle(lipgloss.NormalBorder()).
				BorderLeft(true).
				BorderForeground(theme.Purple).
				PaddingLeft(1).
				Render(rowContent)
		}

		items = append(items, item)
	}

	list := lipgloss.JoinVertical(lipgloss.Left, items...)
	listH := lipgloss.Height(list)
	if listH < h {
		pad := lipgloss.NewStyle().Background(theme.Background).Width(w).Height(h - listH).Render("")
		list = lipgloss.JoinVertical(lipgloss.Left, list, pad)
	}

	border := lipgloss.NewStyle().
		Background(theme.Background).
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true).
		BorderForeground(theme.BorderColor)

	return border.Render(list)
}

func (m Model) renderDetail(w, h int) string {
	if len(m.Profiles) == 0 {
		return lipgloss.NewStyle().
			Background(theme.Background).
			Width(w).
			Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No profiles")
	}

	p := m.Profiles[m.SelectedIndex]
	bg := lipgloss.NewStyle().Background(theme.Background)

	var sections []string

	// Profile section
	sections = append(sections, m.renderDetailSection("PROFILE", w, []detailRow{
		{Key: "Name", Val: theme.DetailValAccentStyle.Render(p.Name)},
		{Key: "Keywords", Val: theme.DetailValStyle.Render(strings.Join(p.Keywords, ", "))},
	}))

	// Filters section
	excludeTraps := theme.DetailValYesStyle.Render("yes")
	if !p.ExcludeTraps {
		excludeTraps = theme.DetailValNoStyle.Render("no")
	}
	requireEmail := theme.DetailValNoStyle.Render("no")
	if p.MustHaveEmail {
		requireEmail = theme.DetailValYesStyle.Render("yes")
	}

	sections = append(sections, m.renderDetailSection("FILTERS", w, []detailRow{
		{Key: "Score range", Val: bg.Foreground(theme.Text).Render(fmt.Sprintf("%d - %d", p.MinScore, p.MaxScore))},
		{Key: "Exclude traps", Val: excludeTraps},
		{Key: "Require email", Val: requireEmail},
		{Key: "Seniority", Val: bg.Foreground(theme.Text).Render(strings.Join(p.SeniorityLevels, ", "))},
	}))

	// Scoring Weights section
	weights := []barRow{
		{Label: "Tech Match", Pct: p.ScoringWeights.TechMatch},
		{Label: "Seniority", Pct: p.ScoringWeights.SeniorityMatch},
		{Label: "Location", Pct: p.ScoringWeights.LocationMatch},
		{Label: "Company", Pct: p.ScoringWeights.CompanyMatch},
		{Label: "Salary", Pct: p.ScoringWeights.SalaryMatch},
		{Label: "Remote", Pct: p.ScoringWeights.RemoteMatch},
	}
	sections = append(sections, m.renderScoringSection("SCORING WEIGHTS", w, weights))

	detail := lipgloss.JoinVertical(lipgloss.Left, sections...)
	detailH := lipgloss.Height(detail)
	if detailH < h {
		pad := bg.Width(w).Height(h - detailH).Render("")
		detail = lipgloss.JoinVertical(lipgloss.Left, detail, pad)
	}

	return lipgloss.NewStyle().
		Background(theme.Background).
		Width(w).
		Padding(1, 2).
		Render(detail)
}

type detailRow struct {
	Key string
	Val string
}

func (m Model) renderDetailSection(title string, w int, rows []detailRow) string {
	bg := lipgloss.NewStyle().Background(theme.Background)

	titleStr := theme.DetailSectionTitleStyle.Render(title)
	border := bg.Foreground(theme.BorderColor).Width(w - 6).Render(strings.Repeat("-", w-6))

	var lines []string
	lines = append(lines, titleStr, border)

	keyW := 16
	for _, r := range rows {
		key := theme.DetailKeyStyle.Width(keyW).Render(r.Key)
		lines = append(lines, bg.Render(key+"  "+r.Val))
	}
	lines = append(lines, bg.Render(""))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

type barRow struct {
	Label string
	Pct   int
}

func (m Model) renderScoringSection(title string, w int, bars []barRow) string {
	bg := lipgloss.NewStyle().Background(theme.Background)

	titleStr := theme.DetailSectionTitleStyle.Render(title)
	border := bg.Foreground(theme.BorderColor).Width(w - 6).Render(strings.Repeat("-", w-6))

	var lines []string
	lines = append(lines, titleStr, border)

	labelW := 14
	barW := w - labelW - 14
	if barW < 10 {
		barW = 10
	}

	for _, b := range bars {
		label := theme.BarLabelStyle.Width(labelW).Align(lipgloss.Right).Render(b.Label)
		filled := barW * b.Pct / 100
		if filled < 0 {
			filled = 0
		}
		empty := barW - filled
		if empty < 0 {
			empty = 0
		}
		bar := theme.BarFillStyle.Render(strings.Repeat("█", filled)) +
			theme.BarTrackStyle.Render(strings.Repeat("░", empty))
		pct := theme.BarPctStyle.Width(5).Align(lipgloss.Right).Render(fmt.Sprintf("%d%%", b.Pct))

		lines = append(lines, bg.Render(label+"  "+bar+" "+pct))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func summarizeKeywords(kw []string, maxW int) string {
	if len(kw) == 0 {
		return "no keywords"
	}
	result := ""
	shown := 0
	for _, k := range kw {
		candidate := result
		if shown > 0 {
			candidate += ", "
		}
		candidate += k
		if len(candidate) > maxW-6 && shown > 0 {
			return result + fmt.Sprintf(", +%d", len(kw)-shown)
		}
		result = candidate
		shown++
	}
	return result
}
