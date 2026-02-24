package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "↓":
			if len(m.jobs) > 0 {
				m.selectedIndex = min(m.selectedIndex+1, len(m.jobs)-1)
				m.viewState = JobList
			}
		case "k", "↑":
			if len(m.jobs) > 0 {
				m.selectedIndex = max(m.selectedIndex-1, 0)
				m.viewState = JobList
			}
		case "s":
			m.viewState = Scraping
		case "f":
			m.viewState = Filter
		case "p":
			m.viewState = Profiles
		case "m":
			m.viewState = Emails
		case "a":
		case "?":
			m.viewState = Help
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}
