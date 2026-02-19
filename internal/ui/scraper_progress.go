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

// ScraperProgressView shows scraping progress with real-time updates
type ScraperProgressView struct {
	progress      progress.Model
	spinner       spinner.Model
	totalSources  int
	currentJob    int
	totalJobs     int
	currentSource string
	elapsedTime   time.Duration
	startTime     time.Time
	width         int
	height        int
	status        string
	results       []string // Recent results for display
}

// NewScraperProgressView creates a new progress view
func NewScraperProgressView() *ScraperProgressView {
	p := progress.New(progress.WithDefaultGradient())
	s := spinner.New(spinner.WithSpinner(spinner.Dot))
	s.Style = lipgloss.NewStyle().Foreground(Colors.Accent)

	return &ScraperProgressView{
		progress:  p,
		spinner:   s,
		status:    "Initializing...",
		startTime: time.Now(),
		results:   make([]string, 0, 5),
	}
}

// SetSize updates the component size
func (p *ScraperProgressView) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.progress.Width = width - 10
}

// Update handles messages
func (p *ScraperProgressView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ScraperProgressMsg:
		p.updateProgress(msg.Progress)
		return nil
	case ScraperJobMsg:
		p.addResult(msg.Job)
		return nil
	case ScraperCompleteMsg:
		p.status = "Complete"
		return nil
	case ScraperErrorMsg:
		p.status = fmt.Sprintf("Error: %v", msg.Error)
		return nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		p.spinner, cmd = p.spinner.Update(msg)
		return cmd
	case progress.FrameMsg:
		progressModel, cmd := p.progress.Update(msg)
		p.progress = progressModel.(progress.Model)
		return cmd
	}

	// Handle spinner and progress updates
	var cmds []tea.Cmd
	_, spinnerCmd := p.spinner.Update(nil)
	if spinnerCmd != nil {
		cmds = append(cmds, spinnerCmd)
	}

	return tea.Batch(cmds...)
}

// View renders the progress view
func (p *ScraperProgressView) View(width, height int) string {
	p.SetSize(width, height)

	header := p.renderHeader()
	progressBar := p.renderProgress()
	status := p.renderStatus()
	results := p.renderResults()
	help := p.renderHelp()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		progressBar,
		status,
		results,
		help,
	)
}

func (p *ScraperProgressView) renderHeader() string {
	title := Styles.Title.Render("Scraping Jobs")
	subtitle := Styles.Subtitle.Render("Searching multiple job sources...")

	return Styles.Card.Width(p.width).Render(
		lipgloss.JoinVertical(lipgloss.Left, title, subtitle),
	)
}

func (p *ScraperProgressView) renderProgress() string {
	return Styles.Content.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			p.progress.View(),
			fmt.Sprintf("Progress: %d/%d sources", p.currentJob, p.totalJobs),
		),
	)
}

func (p *ScraperProgressView) renderStatus() string {
	elapsed := time.Since(p.startTime).Round(time.Second)
	statusText := fmt.Sprintf("%s %s (%s)", p.spinner.View(), p.status, elapsed)

	if p.currentSource != "" {
		statusText += fmt.Sprintf("\nCurrent: %s", p.currentSource)
	}

	return Styles.StatusBar.Width(p.width).Render(
		Styles.StatusText.Render(statusText),
	)
}

func (p *ScraperProgressView) renderResults() string {
	if len(p.results) == 0 {
		return Styles.MutedText.Render("No results yet...")
	}

	var resultText strings.Builder
	resultText.WriteString("Recent results:\n")

	// Show last 3 results
	start := 0
	if len(p.results) > 3 {
		start = len(p.results) - 3
	}

	for i := start; i < len(p.results); i++ {
		resultText.WriteString(fmt.Sprintf("• %s\n", p.results[i]))
	}

	return Styles.Content.Height(5).Render(resultText.String())
}

func (p *ScraperProgressView) renderHelp() string {
	return Styles.StatusBar.Width(p.width).Render(
		Styles.StatusText.Render("Press Ctrl+C to cancel • Esc to return when complete"),
	)
}

func (p *ScraperProgressView) updateProgress(progress scraper.ScraperProgress) {
	p.currentSource = progress.Source
	p.currentJob = progress.CurrentSource
	p.totalJobs = progress.TotalSources
	p.status = progress.Status
	p.elapsedTime = progress.ElapsedTime
}

func (p *ScraperProgressView) addResult(job job.Job) {
	result := fmt.Sprintf("%s @ %s (Score: %d)", job.Title, job.Company, job.Score)

	// Keep only last 5 results
	p.results = append(p.results, result)
	if len(p.results) > 5 {
		p.results = p.results[1:]
	}
}

// Message types for scraper progress
type ScraperProgressMsg struct {
	Progress scraper.ScraperProgress
}

type ScraperJobMsg struct {
	Job job.Job
}

type ScraperCompleteMsg struct{}

type ScraperErrorMsg struct {
	Error error
}

type ScraperCancelMsg struct{}
