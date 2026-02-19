package ui

import (
	"fmt"
	"strings"
)

func (m Model) viewFilters() string {
	return Styles.Border.Width(m.width - 2).Render(
		fmt.Sprintf("Filter by Keywords:\n\n%s", m.filterInput.View()),
	)
}

func (m Model) viewProfiles() string {
	var rows []string
	for _, p := range m.profiles {
		cur, style := "  ", Styles.Text
		if p.ID == m.activeProfile.ID {
			cur, style = "✓ ", Styles.SelectedItem
		}
		rows = append(rows, style.Render(fmt.Sprintf("%s%s (%s)", cur, p.Name, p.ID)))
	}
	return Styles.Border.Width(m.width - 2).Render(strings.Join(rows, "\n"))
}

func (m Model) viewHelp() string {
	help := `
Sprayer - AI-Powered Job Automation

Keybindings:
  s       Scrape jobs (HN, RemoteOK, LinkedIn, etc)
  f       Filter jobs by keyword
  p       Switch profile (persona)
  a       Apply (Generate email draft + attach CV)
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
	return Styles.Border.Width(m.width - 2).Render(help)
}

func truncate(s string, w int) string {
	if len(s) > w {
		return s[:w-3] + "..."
	}
	return s
}
