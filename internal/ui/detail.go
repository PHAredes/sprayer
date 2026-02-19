package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"sprayer/internal/job"
)

type Detail struct {
	Job      job.Job
	Viewport viewport.Model
	Width    int
	Height   int
}

func (d *Detail) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	d.Viewport, cmd = d.Viewport.Update(msg)
	return cmd
}

func (d Detail) View() string {
	var b strings.Builder
	b.WriteString(Styles.Title.Render(d.Job.Title) + "\n")
	b.WriteString(Styles.MutedText.Render(fmt.Sprintf("%s @ %s", d.Job.Company, d.Job.Location)) + "\n\n")
	b.WriteString(fmt.Sprintf("Score: %d • Posted: %s • Source: %s\n",
		d.Job.Score, d.Job.PostedDate.Format("2006-01-02"), d.Job.Source))
	if d.Job.Salary != "" {
		b.WriteString(fmt.Sprintf("Salary: %s\n", d.Job.Salary))
	}
	b.WriteString(fmt.Sprintf("URL: %s\n\n", d.Job.URL))

	if d.Job.HasTraps {
		b.WriteString(Styles.ErrorText.Render(fmt.Sprintf("⚠️  ANTI-AI TRAPS DETECTED: %v", d.Job.Traps)) + "\n\n")
		b.WriteString(Styles.MutedText.Render("(Description below has been sanitized)") + "\n\n")
	}
	b.WriteString(d.Job.Description)

	d.Viewport.SetContent(b.String())
	return Styles.Border.Width(d.Width - 2).Render(d.Viewport.View())
}
