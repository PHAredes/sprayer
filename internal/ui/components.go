package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/internal/job"
	"sprayer/internal/scraper"
)

// Header represents the application header
type Header struct {
	width int
}

func (h *Header) SetWidth(width int) {
	h.width = width
}

func (h *Header) View() string {
	if h.width == 0 {
		return ""
	}

	title := lipgloss.NewStyle().
		Foreground(Colors.Primary).
		Bold(true).
		Render("Sprayer")

	subtitle := lipgloss.NewStyle().
		Foreground(Colors.Muted).
		Render("Job Application Tool")

	content := lipgloss.JoinHorizontal(
		lipgloss.Center,
		title,
		" ",
		subtitle,
	)

	return Styles.Header.
		Width(h.width).
		Render(content)
}

// Footer represents the application footer with keybinds
type Footer struct {
	width int
}

func (f *Footer) SetWidth(width int) {
	f.width = width
}

func (f *Footer) View(mode AppMode) string {
	if f.width == 0 {
		return ""
	}

	var keybinds string
	switch mode {
	case ModeJobs:
		keybinds = "s: scrape | r: refresh | ↑↓: navigate | enter: select | ?: help | q: quit"
	case ModeScraping:
		keybinds = "esc: cancel | ?: help | q: quit"
	case ModeHelp:
		keybinds = "esc/q: back to jobs"
	default:
		keybinds = "?: help | q: quit"
	}

	return Styles.StatusBar.
		Width(f.width).
		Render(keybinds)
}

// SimpleJobList represents the job listing component
type SimpleJobList struct {
	jobs     []job.Job
	selected int
	width    int
	height   int
}

func (j *SimpleJobList) SetJobs(jobs []job.Job) {
	j.jobs = jobs
}

func (j *SimpleJobList) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if j.selected > 0 {
				j.selected--
			}
		case "down", "j":
			if j.selected < len(j.jobs)-1 {
				j.selected++
			}
		}
	}
	return nil
}

func (j *SimpleJobList) View() string {
	if j.width == 0 || j.height == 0 {
		return "Loading jobs..."
	}

	if len(j.jobs) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(Colors.Muted).
			Align(lipgloss.Center).
			Width(j.width).
			Height(j.height - 1).
			Render("No jobs found. Press 's' to scrape for jobs.")
		return emptyStyle
	}

	var lines []string
	for i, job := range j.jobs {
		style := Styles.ListItem
		if i == j.selected {
			style = Styles.SelectedItem
		}

		trapIndicator := ""
		if job.HasTraps {
			trapIndicator = " [!]"
		}

		line := fmt.Sprintf("[%d]%s %s @ %s", job.Score, trapIndicator, job.Title, job.Company)
		lines = append(lines, style.Render(line))
	}

	return lipgloss.NewStyle().
		Height(j.height - 1).
		Render(strings.Join(lines, "\n"))
}

func (j *SimpleJobList) SetSize(width, height int) {
	j.width = width
	j.height = height
}

func (j *SimpleJobList) SelectedJob() *job.Job {
	if j.selected >= 0 && j.selected < len(j.jobs) {
		return &j.jobs[j.selected]
	}
	return nil
}

func (j *SimpleJobList) ToggleSort() {
	// Simple implementation - just reverse the list for now
	for i, k := 0, len(j.jobs)-1; i < k; i, k = i+1, k-1 {
		j.jobs[i], j.jobs[k] = j.jobs[k], j.jobs[i]
	}
}

// ScraperView represents the scraping progress view
type ScraperView struct {
	width    int
	height   int
	progress scraper.ScraperProgress
	error    error
}

func (s *ScraperView) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (s *ScraperView) UpdateProgress(progress scraper.ScraperProgress) tea.Cmd {
	s.progress = progress
	return nil
}

func (s *ScraperView) SetError(err error) {
	s.error = err
}

func (s *ScraperView) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s *ScraperView) View() string {
	if s.width == 0 || s.height == 0 {
		return ""
	}

	content := lipgloss.NewStyle().
		Foreground(Colors.Accent).
		Bold(true).
		Render("🔍 Scraping jobs...")

	if s.progress.Source != "" {
		progressText := fmt.Sprintf("\nSource: %s", s.progress.Source)
		if s.progress.JobsFound > 0 {
			progressText += fmt.Sprintf("\nJobs found: %d", s.progress.JobsFound)
		}
		if s.progress.TotalSources > 0 {
			progressText += fmt.Sprintf("\nProgress: %d/%d sources",
				s.progress.CurrentSource, s.progress.TotalSources)
		}
		if s.progress.ElapsedTime > 0 {
			progressText += fmt.Sprintf("\nElapsed: %s",
				s.progress.ElapsedTime.Round(time.Second))
		}
		if s.progress.Status != "" {
			progressText += fmt.Sprintf("\nStatus: %s", s.progress.Status)
		}

		content += lipgloss.NewStyle().
			Foreground(Colors.Text).
			MarginTop(1).
			Render(progressText)
	}

	if s.error != nil {
		errorText := lipgloss.NewStyle().
			Foreground(Colors.Error).
			Render(fmt.Sprintf("\nError: %v", s.error))
		content += errorText
	}

	return Styles.Scraping.
		Width(s.width).
		Height(s.height).
		Render(content)
}
