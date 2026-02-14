package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"sprayer/internal/job"
	"sprayer/internal/profile"
)

func (m Model) viewJobList() string {
	header := styleTitle.Render(fmt.Sprintf("Sprayer • %s (%d jobs)", m.activeProfile.Name, len(m.filteredJobs)))
	
	// List content
	var rows []string
	start, end := m.visibleRange()
	
	for i := start; i < end; i++ {
		j := m.filteredJobs[i]
		cursor := "  "
		style := styleNormal
		if i == m.cursor {
			cursor = "> "
			style = styleSelected
		}
		
		line := fmt.Sprintf("%s%s • %s", cursor, j.Title, j.Company)
		rows = append(rows, style.Render(truncate(line, m.width-4)))
	}
	
	content := styleBox.
		Width(m.width - 2).
		Height(m.height - 4). // Room for header + footer
		Render(strings.Join(rows, "\n"))

	// Footer
	keys := "s: scrape • f: filter • p: profiles • /: search • enter: view • q: quit"
	if m.err != nil {
		keys = styleError.Render(fmt.Sprintf("Error: %v", m.err))
	} else if m.statusMsg != "" {
		keys = styleSuccess.Render(m.statusMsg)
	}
	
	footer := styleFooter.Width(m.width).Render(keys)

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

func (m Model) viewJobDetail() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		styleTitle.Render("Job Details (Esc to back, 'a' to apply)"),
		styleBox.Width(m.width-2).Render(m.viewport.View()),
		styleFooter.Width(m.width).Render("j/k: scroll • a: apply • esc: back"),
	)
}

func (m Model) viewFilters() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		styleTitle.Render("Filters"),
		styleBox.Width(m.width-2).Render(
			fmt.Sprintf("Filter by Keywords:\n\n%s", m.filterInput.View()),
		),
		styleFooter.Width(m.width).Render("Enter: apply • Esc: cancel"),
	)
}

func (m Model) viewProfiles() string {
	var rows []string
	for _, p := range m.profiles {
		cursor := "  "
		style := styleNormal
		if p.ID == m.activeProfile.ID {
			cursor = "✓ "
			style = styleSelected
		}
		rows = append(rows, style.Render(fmt.Sprintf("%s%s (%s)", cursor, p.Name, p.ID)))
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		styleTitle.Render("Select Profile"),
		styleBox.Width(m.width-2).Render(strings.Join(rows, "\n")),
		styleFooter.Width(m.width).Render("↑/↓: select • Enter: confirm"),
	)
}

func (m Model) viewHelp() string {
	help := `
Sprayer - AI-Powered Job Automation

Keybindings:
  s       Scrape jobs (HN, RemoteOK, LinkedIn, etc)
  f       Filter jobs by keyword
  p       Switch profile (persona)
  a       Apply (Generate email draft + attach CV)
  /       Search within list
  j/k     Navigate list
  Enter   View details
  q       Quit

Configuration:
  ~/.sprayer/         Data directory (SQLite DB)
  ~/.sprayer/prompts  LLM prompts
  
Env Vars:
  SPRAYER_LLM_KEY     API Key
  SPRAYER_LLM_MODEL   Model (default: kimi-k2)
`
	return lipgloss.JoinVertical(lipgloss.Left,
		styleTitle.Render("Help"),
		styleBox.Width(m.width-2).Render(help),
		styleFooter.Width(m.width).Render("Press any key to return"),
	)
}

func (m Model) visibleRange() (int, int) {
	height := m.height - 6 // Header + footer + box padding
	if height <= 0 {
		height = 1
	}
	
	start := m.cursor - (height / 2)
	if start < 0 {
		start = 0
	}
	
	end := start + height
	if end > len(m.filteredJobs) {
		end = len(m.filteredJobs)
		start = end - height
		if start < 0 {
			start = 0
		}
	}
	return start, end
}

func renderJobDetail(j job.Job, p profile.Profile) string {
	var b strings.Builder
	
	b.WriteString(styleTitle.Render(j.Title) + "\n")
	b.WriteString(styleDim.Render(fmt.Sprintf("%s @ %s", j.Company, j.Location)) + "\n\n")
	
	// Stats
	b.WriteString(fmt.Sprintf("Score: %d • Posted: %s • Source: %s\n", 
		j.Score, j.PostedDate.Format("2006-01-02"), j.Source))
	if j.Salary != "" {
		b.WriteString(fmt.Sprintf("Salary: %s\n", j.Salary))
	}
	b.WriteString(fmt.Sprintf("URL: %s\n\n", j.URL))
	
	b.WriteString(j.Description)
	
	return b.String()
}

func truncate(s string, w int) string {
	if len(s) > w {
		return s[:w-3] + "..."
	}
	return s
}
