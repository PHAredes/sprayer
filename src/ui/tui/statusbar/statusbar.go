package statusbar

import (
	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/theme"
)

// Binding represents a single key-label pair in the status bar.
type Binding struct {
	Key   string
	Label string
}

// Render builds a complete status bar row with key bindings on the left
// and an optional right-aligned section.
func Render(width int, bindings []Binding, right string) string {
	kbd := lipgloss.NewStyle().Background(theme.Surface).Foreground(theme.Cyan)
	sp := lipgloss.NewStyle().Background(theme.Surface).Foreground(theme.Subtle).Render(" ")
	label := lipgloss.NewStyle().Background(theme.Surface).Foreground(theme.Subtle)
	sep := theme.SepStyle.Render(" | ")

	line := ""
	for i, b := range bindings {
		if i > 0 {
			line += sep
		}
		line += kbd.Render(b.Key) + sp + label.Render(b.Label)
	}

	if right != "" {
		rightStyle := lipgloss.NewStyle().Background(theme.Surface).Foreground(theme.Subtle)
		rightRendered := rightStyle.Render(right)
		leftW := lipgloss.Width(line)
		rightW := lipgloss.Width(rightRendered)
		gap := width - leftW - rightW - 4
		if gap < 1 {
			gap = 1
		}
		pad := lipgloss.NewStyle().Background(theme.Surface).Width(gap).Render("")
		line += pad + rightRendered
	}

	return lipgloss.NewStyle().
		Background(theme.Surface).
		Width(width).
		PaddingLeft(2).
		PaddingRight(2).
		Render(line)
}
