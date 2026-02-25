package help

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/theme"
)

type section struct {
	Title string
	Rows  []row
}

type row struct {
	Key  string
	Desc string
}

// Model holds only dimensions — help is stateless.
type Model struct {
	Width  int
	Height int
}

func New(w, h int) Model {
	return Model{Width: w, Height: h}
}

func (m Model) View() string {
	sections := []section{
		{
			Title: "Job List",
			Rows: []row{
				{"s", "Scrape for new jobs"},
				{"f", "Open filters"},
				{"p", "Manage profiles"},
				{"m", "View email drafts"},
				{"up/dn/jk", "Navigate"},
				{"enter", "Job detail"},
				{"a", "Apply to job"},
			},
		},
		{
			Title: "Job Detail",
			Rows: []row{
				{"esc", "Back to list"},
				{"e", "Generate email"},
				{"a", "Gen application"},
			},
		},
		{
			Title: "Emails",
			Rows: []row{
				{"e", "New email for job"},
				{"enter", "Open compose"},
				{"ctrl+s", "Send via pop"},
				{"r", "Regenerate (mods)"},
				{"d", "Delete draft"},
			},
		},
		{
			Title: "Global",
			Rows: []row{
				{"?", "This help"},
				{"esc", "Back / cancel"},
				{"ctrl+c", "Quit"},
			},
		},
	}

	cols := m.renderColumns(sections)
	return cols
}

func (m Model) renderColumns(sections []section) string {
	availH := m.Height - 2 // minus top/statusbar
	if availH < 1 {
		availH = 1
	}

	numCols := 2
	colW := (m.Width - 6) / numCols
	if colW < 30 {
		numCols = 1
		colW = m.Width - 4
	}

	columns := make([]string, numCols)
	sIdx := 0
	for c := 0; c < numCols && sIdx < len(sections); c++ {
		var parts []string
		// distribute roughly evenly
		end := sIdx + (len(sections)-sIdx+numCols-c-1)/(numCols-c)
		for ; sIdx < end && sIdx < len(sections); sIdx++ {
			parts = append(parts, m.renderSection(sections[sIdx], colW))
		}
		columns[c] = lipgloss.JoinVertical(lipgloss.Left, parts...)
	}

	joined := lipgloss.JoinHorizontal(lipgloss.Top, columns...)

	return lipgloss.Place(
		m.Width, availH,
		lipgloss.Center, lipgloss.Top,
		joined,
		lipgloss.WithWhitespaceBackground(theme.Background),
	)
}

func (m Model) renderSection(s section, colW int) string {
	bg := lipgloss.NewStyle().Background(theme.Background)

	title := theme.HelpSectionTitleStyle.Render(s.Title)
	border := bg.Foreground(theme.BorderColor).Width(colW - 4).Render(strings.Repeat("-", colW-4))
	header := lipgloss.JoinVertical(lipgloss.Left, title, border)

	var rows []string
	for _, r := range s.Rows {
		keyW := 12
		key := theme.HelpKeyStyle.Width(keyW).Align(lipgloss.Center).Render(r.Key)
		desc := theme.HelpDescStyle.Render(" " + r.Desc)
		rows = append(rows, bg.Render(key+desc))
	}

	block := lipgloss.JoinVertical(lipgloss.Left, append([]string{header, bg.Render("")}, rows...)...)

	return lipgloss.NewStyle().
		Background(theme.Background).
		Width(colW).
		Padding(1, 2).
		Render(block)
}
