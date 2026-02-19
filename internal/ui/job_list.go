package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"sprayer/internal/job"
)

// JobList represents a CHARM-style job list component
type JobList struct {
	jobs         []job.Job
	filteredJobs []job.Job
	cursor       int
	width        int
	height       int
	sortMode     SortMode
	selectedJob  *job.Job
}

type SortMode int

const (
	SortByScore SortMode = iota
	SortByDate
	SortByTitle
	SortByCompany
)

// NewJobList creates a new CHARM-style job list
func NewJobList(jobs []job.Job) *JobList {
	return &JobList{
		jobs:         jobs,
		filteredJobs: jobs,
		cursor:       0,
		sortMode:     SortByScore,
	}
}

// SetJobs updates the job list
func (l *JobList) SetJobs(jobs []job.Job) {
	l.jobs = jobs
	l.filteredJobs = jobs
	l.cursor = 0
	l.applySorting()
}

// SetSize updates the component size
func (l *JobList) SetSize(width, height int) {
	l.width = width
	l.height = height
}

// SelectedJob returns the currently selected job
func (l *JobList) SelectedJob() *job.Job {
	if l.cursor >= 0 && l.cursor < len(l.filteredJobs) {
		return &l.filteredJobs[l.cursor]
	}
	return nil
}

// ToggleSort cycles through sort modes
func (l *JobList) ToggleSort() {
	l.sortMode = (l.sortMode + 1) % 4
	l.applySorting()
}

// Update handles messages
func (l *JobList) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Up):
			if l.cursor > 0 {
				l.cursor--
			}
		case key.Matches(msg, Keys.Down):
			if l.cursor < len(l.filteredJobs)-1 {
				l.cursor++
			}
		case key.Matches(msg, Keys.Enter):
			if l.SelectedJob() != nil {
				return func() tea.Msg {
					return JobSelectedMsg{Job: l.SelectedJob()}
				}
			}
		}
	}
	return nil
}

// View renders the component
func (l *JobList) View(width, height int) string {
	l.width = width
	l.height = height

	if len(l.filteredJobs) == 0 {
		return l.renderEmptyState()
	}

	header := l.renderHeader()
	list := l.renderList()

	return lipgloss.JoinVertical(lipgloss.Left, header, list)
}

func (l *JobList) renderHeader() string {
	sortIndicator := ""
	switch l.sortMode {
	case SortByScore:
		sortIndicator = "📊 Score"
	case SortByDate:
		sortIndicator = "📅 Date"
	case SortByTitle:
		sortIndicator = "📝 Title"
	case SortByCompany:
		sortIndicator = "🏢 Company"
	}

	left := Styles.HeaderText.Render(fmt.Sprintf("%d jobs", len(l.filteredJobs)))
	center := Styles.Title.Render("Job Listings")
	right := Styles.HeaderText.Render(sortIndicator)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		Styles.Header.Width(l.width/3).Render(left),
		Styles.Header.Width(l.width/3).Align(lipgloss.Center).Render(center),
		Styles.Header.Width(l.width/3).Align(lipgloss.Right).Render(right),
	)
}

func (l *JobList) renderList() string {
	visibleHeight := l.height - 3 // Account for header and status
	if visibleHeight <= 0 {
		visibleHeight = 1
	}

	start := l.cursor - (visibleHeight / 2)
	if start < 0 {
		start = 0
	}
	end := start + visibleHeight
	if end > len(l.filteredJobs) {
		end = len(l.filteredJobs)
		start = end - visibleHeight
		if start < 0 {
			start = 0
		}
	}

	var rows []string
	for i := start; i < end; i++ {
		job := l.filteredJobs[i]
		row := l.renderJobRow(job, i == l.cursor)
		rows = append(rows, row)
	}

	// Fill empty space
	for i := len(rows); i < visibleHeight; i++ {
		rows = append(rows, strings.Repeat(" ", l.width))
	}

	return strings.Join(rows, "\n")
}

func (l *JobList) renderJobRow(job job.Job, selected bool) string {
	// Truncate title and company to fit
	maxTitleLen := l.width/3 - 2
	maxCompanyLen := l.width/3 - 2
	maxLocationLen := l.width/3 - 8 // Account for score

	title := truncateJobString(job.Title, maxTitleLen)
	company := truncateJobString(job.Company, maxCompanyLen)
	location := truncateJobString(job.Location, maxLocationLen)
	score := fmt.Sprintf("%3d", job.Score)

	// Create columns
	titleCol := fmt.Sprintf("%-*s", maxTitleLen, title)
	companyCol := fmt.Sprintf("%-*s", maxCompanyLen, company)
	locationCol := fmt.Sprintf("%-*s", maxLocationLen, location)
	scoreCol := fmt.Sprintf("%3s", score)

	row := fmt.Sprintf("%s %s %s %s", titleCol, companyCol, locationCol, scoreCol)

	if selected {
		return Styles.SelectedItem.Width(l.width).Render(row)
	}

	// Score-based coloring
	if job.Score >= 80 {
		return Styles.SuccessText.Width(l.width).Render(row)
	} else if job.Score >= 60 {
		return Styles.Text.Width(l.width).Render(row)
	} else {
		return Styles.MutedText.Width(l.width).Render(row)
	}
}

func (l *JobList) renderEmptyState() string {
	emptyMsg := "No jobs found matching your criteria.\n\nTry adjusting your filters or scrape for new jobs."

	return lipgloss.Place(l.width, l.height-1,
		lipgloss.Center, lipgloss.Center,
		Styles.MutedText.Render(emptyMsg),
	)
}

func (l *JobList) applySorting() {
	// Sort based on current mode
	switch l.sortMode {
	case SortByScore:
		// Already sorted by score (default)
	case SortByDate:
		// Sort by date (newest first)
		// Implementation would sort l.filteredJobs by PostedDate
	case SortByTitle:
		// Sort by title alphabetically
		// Implementation would sort l.filteredJobs by Title
	case SortByCompany:
		// Sort by company alphabetically
		// Implementation would sort l.filteredJobs by Company
	}
}

// MoveCursor moves the cursor up or down
func (l *JobList) MoveCursor(up bool) {
	if up && l.cursor > 0 {
		l.cursor--
	} else if !up && l.cursor < len(l.filteredJobs)-1 {
		l.cursor++
	}
}

// JobSelectedMsg is sent when a job is selected
type JobSelectedMsg struct {
	Job *job.Job
}

// Helper function to truncate strings
func truncateJobString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
