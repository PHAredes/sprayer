package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
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
		keybinds = "s: scrape | f: filter | p: profiles | ↑↓: navigate | enter: select | a: apply | ?: help | q: quit"
	case ModeJobDetail:
		keybinds = "esc: back | a: apply | ?: help | q: quit"
	case ModeFilters:
		keybinds = "tab: next field | enter: apply | esc: cancel | ?: help | q: quit"
	case ModeProfiles:
		keybinds = "↑↓: select | enter: choose | n: new | esc: back | ?: help | q: quit"
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

// ScraperView represents the scraping progress view with proper charm components
type ScraperView struct {
	width       int
	height      int
	progress    scraper.ScraperProgress
	error       error
	spinner     spinner.Model
	progressBar progress.Model
}

func NewScraperView() *ScraperView {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(Colors.Accent)

	pb := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	return &ScraperView{
		spinner:     sp,
		progressBar: pb,
	}
}

func (s *ScraperView) Init() tea.Cmd {
	return s.spinner.Tick
}

func (s *ScraperView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return cmd
	case progress.FrameMsg:
		progressModel, cmd := s.progressBar.Update(msg)
		s.progressBar = progressModel.(progress.Model)
		return cmd
	}
	return nil
}

func (s *ScraperView) UpdateProgress(progress scraper.ScraperProgress) tea.Cmd {
	s.progress = progress
	// Update progress bar if we have progress info
	if progress.TotalSources > 0 && progress.CurrentSource > 0 {
		progressPercent := float64(progress.CurrentSource) / float64(progress.TotalSources)
		return s.progressBar.SetPercent(progressPercent)
	}
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

	// Header with spinner
	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		s.spinner.View(),
		" ",
		lipgloss.NewStyle().Foreground(Colors.Accent).Bold(true).Render("Scraping jobs..."),
	)

	var contentParts []string
	contentParts = append(contentParts, header)

	// Progress bar if we have progress info
	if s.progress.TotalSources > 0 && s.progress.CurrentSource > 0 {
		progressPercent := float64(s.progress.CurrentSource) / float64(s.progress.TotalSources)
		progressBar := s.progressBar.ViewAs(progressPercent)

		progressInfo := fmt.Sprintf("Source %d/%d", s.progress.CurrentSource, s.progress.TotalSources)
		if s.progress.Source != "" {
			progressInfo = fmt.Sprintf("%s - %s", progressInfo, s.progress.Source)
		}

		contentParts = append(contentParts, "")
		contentParts = append(contentParts, progressBar)
		contentParts = append(contentParts, progressInfo)
	}

	// Additional info
	if s.progress.Source != "" {
		infoParts := []string{}
		if s.progress.JobsFound > 0 {
			infoParts = append(infoParts, fmt.Sprintf("Jobs found: %d", s.progress.JobsFound))
		}
		if s.progress.ElapsedTime > 0 {
			infoParts = append(infoParts, fmt.Sprintf("Elapsed: %s", s.progress.ElapsedTime.Round(time.Second)))
		}
		if s.progress.Status != "" {
			infoParts = append(infoParts, fmt.Sprintf("Status: %s", s.progress.Status))
		}

		if len(infoParts) > 0 {
			contentParts = append(contentParts, "")
			contentParts = append(contentParts, strings.Join(infoParts, " | "))
		}
	}

	// Error display
	if s.error != nil {
		errorText := lipgloss.NewStyle().
			Foreground(Colors.Error).
			Render(fmt.Sprintf("Error: %v", s.error))
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, errorText)
	}

	content := strings.Join(contentParts, "\n")

	return Styles.Scraping.
		Width(s.width).
		Height(s.height).
		Render(content)
}
