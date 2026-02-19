package ui

import (
	"fmt"
	"strings"

	bubblekey "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"sprayer/internal/job"
)

type List struct {
	Jobs   []job.Job
	Cursor int
	Width  int
	Height int
}

func (l *List) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case bubblekey.Matches(msg, Keys.Up):
			if l.Cursor > 0 {
				l.Cursor--
			}
		case bubblekey.Matches(msg, Keys.Down):
			if l.Cursor < len(l.Jobs)-1 {
				l.Cursor++
			}
		}
	}
	return nil
}

func (l List) View() string {
	var rows []string
	h := l.Height - 6
	if h <= 0 {
		h = 1
	}

	start := l.Cursor - (h / 2)
	if start < 0 {
		start = 0
	}
	end := start + h
	if end > len(l.Jobs) {
		end = len(l.Jobs)
		start = end - h
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end; i++ {
		j := l.Jobs[i]
		cur, style := "  ", Styles.Text
		if i == l.Cursor {
			cur, style = "> ", Styles.SelectedItem
		}
		icon := ""
		if j.HasTraps {
			icon = " [!]"
		}
		line := fmt.Sprintf("%s%s • %s%s", cur, j.Title, j.Company, icon)
		rows = append(rows, style.Render(truncate(line, l.Width-4)))
	}

	return Styles.Border.Width(l.Width - 2).Height(l.Height - 4).Render(strings.Join(rows, "\n"))
}
