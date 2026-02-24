package joblist

import (

	"github.com/charmbracelet/lipgloss"
	"sprayer/src/api/job"
	"sprayer/src/ui/tui/theme"
)

type Model struct {
	Jobs          []job.Job
	SelectedIndex int
	Width         int
	Height        int
}

func (m Model) View() string {
	if len(m.Jobs) == 0 {
		return m.renderEmptyState()
	}
	return m.renderJobList()
}

func (m Model) contentHeight() int { return m.Height - 2 }

func (m Model) renderEmptyState() string {
	availH := m.contentHeight()
	bg := lipgloss.NewStyle().Background(theme.Background)

	ascii := theme.EmptyASCIIStyle.Render(
		"┌─────────────────────────┐\n" +
			"│  ░░░░░░░░░░░░░░░░░░░░░  │\n" +
			"│  ░  No jobs cached.  ░  │\n" +
			"│  ░░░░░░░░░░░░░░░░░░░░░  │\n" +
			"└─────────────────────────┘")

	headline := theme.EmptyHeadlineStyle.Render("No jobs found")

	subLines := m.emptySubLines()
	subW := 0
	for _, l := range subLines {
		if w := lipgloss.Width(l); w > subW {
			subW = w
		}
	}
	padded := make([]string, len(subLines))
	for i, l := range subLines {
		padded[i] = bg.Width(subW).Render(l)
	}
	subtext := lipgloss.JoinVertical(lipgloss.Left, padded...)

	sp := bg.Render("")
	blockW := max(lipgloss.Width(ascii), lipgloss.Width(headline), subW)

	wrap := func(s string) string {
		return bg.Width(blockW).Align(lipgloss.Center).Render(s)
	}

	block := lipgloss.JoinVertical(lipgloss.Left,
		wrap(ascii),
		wrap(sp),
		wrap(headline),
		wrap(sp),
		wrap(subtext),
	)

	return lipgloss.Place(
		m.Width, availH,
		lipgloss.Center, lipgloss.Center,
		block,
		lipgloss.WithWhitespaceBackground(theme.Background),
	)
}

func (m Model) emptySubLines() []string {
	prose := lipgloss.NewStyle().Background(theme.Background).Foreground(theme.Subtle)
	kbd := lipgloss.NewStyle().Background(lipgloss.Color("#0d2b33")).Foreground(theme.Cyan).
		PaddingLeft(1).PaddingRight(1)
	dot := prose.Render(" · ")

	return []string{
		prose.Render("Press ") + kbd.Render("s") + prose.Render(" to scrape from all configured sources."),
		prose.Render("Press ") + kbd.Render("f") + prose.Render(" to set filters") + dot + kbd.Render("p") + prose.Render(" to manage profiles."),
		prose.Render("Jobs appear here in real time as they're discovered."),
	}
}

func (m Model) renderJobList() string {
	availH := m.contentHeight()

	var lines []string
	for i, j := range m.Jobs {
		var line string
		if i == m.SelectedIndex {
			line = theme.JobItemSelectedStyle.Width(m.Width).Render(m.formatJobItem(j))
		} else {
			line = theme.JobItemStyle.Width(m.Width).Render(m.formatJobItem(j))
		}
		lines = append(lines, line)
	}
	for len(lines) < availH {
		lines = append(lines, theme.ContentStyle.Width(m.Width).Render(""))
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m Model) formatJobItem(j job.Job) string {
	scoreStr := theme.JobScoreStyle.Render("[" + string(rune('0'+j.Score/10)) + string(rune('0'+j.Score%10)) + "]")
	companyStr := theme.JobCompanyStyle.Render("@ " + j.Company)
	sourceStr := theme.JobSourceStyle.Render("(" + j.Source + ")")
	trapStr := theme.JobTrapsStyle.Render(" [!]")
	traps := ""
	if j.HasTraps {
		traps = trapStr
	}

	availW := m.Width - lipgloss.Width(scoreStr) - lipgloss.Width(companyStr) -
		lipgloss.Width(sourceStr) - lipgloss.Width(traps) - 4
	title := j.Title
	if len(title) > availW && availW > 3 {
		title = title[:availW-3] + "..."
	}
	titleStr := theme.JobItemStyle.Render(title)

	return scoreStr + " " + titleStr + " " + companyStr + " " + sourceStr + traps
}
