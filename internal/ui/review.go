package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/internal/job"
)

type Review struct {
	Job     job.Job
	Input   textarea.Model
	Traps   []string
	Subject string
	Width   int
	Height  int
	CVPath  string
}

func (r *Review) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	r.Input, cmd = r.Input.Update(msg)
	return cmd
}

func (r Review) View() string {
	head := fmt.Sprintf("Review Email to %s @ %s", r.Job.Title, r.Job.Company)
	if len(r.Traps) > 0 {
		head += fmt.Sprintf("\n⚠️  TRAPS: %v", r.Traps)
	}

	title := Styles.Title.Render(head)
	body := Styles.FilterBox.Render(r.Input.View())
	foot := Styles.StatusBar.Width(r.Width).Render(fmt.Sprintf("Attachment: %s • Ctrl+Enter: Send • Esc: Cancel", r.CVPath))

	return lipgloss.JoinVertical(lipgloss.Left, title, body, foot)
}
