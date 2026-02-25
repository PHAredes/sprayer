package scraping

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/theme"
)

// spinnerFrames are the Unicode braille spinner characters.
var spinnerFrames = []string{"◐", "◓", "◑", "◒"}

// TickMsg triggers spinner animation.
type TickMsg time.Time

// CancelMsg signals scraping cancellation.
type CancelMsg struct{}

// Model holds scraping progress state.
type Model struct {
	Current    int
	Total      int
	SourceName string
	Found      int
	StartTime  time.Time
	Done       bool
	Frame      int
	Width      int
	Height     int
}

// New creates a new scraping model.
func New(total int, w, h int) Model {
	return Model{
		Current:    0,
		Total:      total,
		SourceName: "",
		Found:      0,
		StartTime:  time.Now(),
		Done:       false,
		Frame:      0,
		Width:      w,
		Height:     h,
	}
}

// Tick returns a command that sends TickMsg after a delay.
func Tick() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Update handles scraping screen events.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		_ = msg
		m.Frame++
		if !m.Done {
			return m, Tick()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return CancelMsg{} }
		}
	}
	return m, nil
}

// View renders the scraping progress screen.
func (m Model) View() string {
	contentH := m.Height - 2
	if contentH < 1 {
		contentH = 1
	}

	bg := lipgloss.NewStyle().Background(theme.Background)

	// Spinner + title
	spinChar := spinnerFrames[m.Frame%len(spinnerFrames)]
	spinner := theme.SpinnerStyle.Render(spinChar)
	title := lipgloss.NewStyle().Foreground(theme.Bright).Bold(true).Render(" Scraping jobs...")

	titleLine := bg.Render(spinner + title)

	// Progress bar
	barW := m.Width * 60 / 100
	if barW < 20 {
		barW = 20
	}
	if barW > 60 {
		barW = 60
	}

	pct := 0
	if m.Total > 0 {
		pct = m.Current * 100 / m.Total
	}

	filled := barW * pct / 100
	empty := barW - filled
	if empty < 0 {
		empty = 0
	}

	bar := theme.ProgressFillStyle.Render(strings.Repeat("█", filled)) +
		theme.ProgressTrackStyle.Render(strings.Repeat("░", empty))

	// Meta line
	sourceInfo := lipgloss.NewStyle().Foreground(theme.Subtle).Render(
		fmt.Sprintf("Source %d/%d", m.Current, m.Total))
	sourceName := ""
	if m.SourceName != "" {
		sourceName = " · " + lipgloss.NewStyle().Foreground(theme.Cyan).Render(m.SourceName)
	}
	pctStr := theme.ProgressPctStyle.Render(fmt.Sprintf("%d%%", pct))

	metaLeft := sourceInfo + sourceName
	metaLine := bg.Width(barW).Render(metaLeft +
		strings.Repeat(" ", max(1, barW-lipgloss.Width(metaLeft)-lipgloss.Width(pctStr))) +
		pctStr)

	// Elapsed + found
	elapsed := time.Since(m.StartTime).Truncate(time.Second)
	subLine1 := lipgloss.NewStyle().Foreground(theme.Subtle).Render("Elapsed: ") +
		lipgloss.NewStyle().Foreground(theme.Cyan).Render(elapsed.String()) +
		lipgloss.NewStyle().Foreground(theme.Subtle).Render("  ·  Found: ") +
		lipgloss.NewStyle().Foreground(theme.Cyan).Render(fmt.Sprintf("%d", m.Found)) +
		lipgloss.NewStyle().Foreground(theme.Subtle).Render(" jobs so far")

	subLine2 := lipgloss.NewStyle().Foreground(theme.Subtle).Render(
		"Jobs appear in real time as they're discovered.")

	// Center everything
	center := func(s string) string {
		return bg.Width(barW).Align(lipgloss.Center).Render(s)
	}

	block := lipgloss.JoinVertical(lipgloss.Left,
		center(titleLine),
		bg.Render(""),
		center(bar),
		metaLine,
		bg.Render(""),
		center(subLine1),
		center(subLine2),
	)

	return lipgloss.Place(
		m.Width, contentH,
		lipgloss.Center, lipgloss.Center,
		block,
		lipgloss.WithWhitespaceBackground(theme.Background),
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
