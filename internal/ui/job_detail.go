package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"sprayer/internal/job"
)

// JobDetail represents a CHARM-style job detail view
type JobDetail struct {
	job      *job.Job
	viewport viewport.Model
	width    int
	height   int
}

// NewJobDetail creates a new CHARM-style job detail view
func NewJobDetail() *JobDetail {
	vp := viewport.New(80, 24)
	vp.Style = Styles.Detail

	return &JobDetail{
		viewport: vp,
	}
}

// SetJob updates the displayed job
func (d *JobDetail) SetJob(job *job.Job) {
	d.job = job
	d.updateContent()
}

// SetSize updates the component size
func (d *JobDetail) SetSize(width, height int) {
	d.width = width
	d.height = height
	d.viewport.Width = width - 4
	d.viewport.Height = height - 6
}

// Update handles messages
func (d *JobDetail) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	d.viewport, cmd = d.viewport.Update(msg)
	return cmd
}

// View renders the component
func (d *JobDetail) View(width, height int) string {
	d.SetSize(width, height)

	if d.job == nil {
		return d.renderEmpty()
	}

	header := d.renderHeader()
	content := d.viewport.View()
	footer := d.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		footer,
	)
}

func (d *JobDetail) renderHeader() string {
	if d.job == nil {
		return ""
	}

	// Job title and company
	title := Styles.DetailTitle.Render(d.job.Title)
	company := Styles.Subtitle.Render(fmt.Sprintf("🏢 %s", d.job.Company))
	location := Styles.MutedText.Render(fmt.Sprintf("📍 %s", d.job.Location))

	// Score and metadata
	scoreColor := Colors.Success
	if d.job.Score < 60 {
		scoreColor = Colors.Warning
	}
	if d.job.Score < 40 {
		scoreColor = Colors.Error
	}

	score := lipgloss.NewStyle().
		Foreground(scoreColor).
		Bold(true).
		Render(fmt.Sprintf("⭐ Score: %d/100", d.job.Score))

	// Trap indicator
	trapIndicator := ""
	if d.job.HasTraps {
		trapIndicator = Styles.ErrorText.Render("⚠️  Contains application traps")
	}

	header := lipgloss.JoinVertical(lipgloss.Left,
		title,
		company,
		location,
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, score, "  ", trapIndicator),
	)

	return Styles.Card.Width(d.width).Render(header)
}

func (d *JobDetail) renderFooter() string {
	if d.job == nil {
		return ""
	}

	// Action hints
	actions := []string{
		Styles.MutedText.Render("Apply"),
		Styles.MutedText.Render("Back"),
	}

	// Job metadata
	meta := []string{}
	if d.job.Email != "" {
		meta = append(meta, Styles.SuccessText.Render("✉️ Has contact email"))
	}
	if d.job.Salary != "" {
		meta = append(meta, Styles.Text.Render(fmt.Sprintf("💰 %s", d.job.Salary)))
	}
	if d.job.Source != "" {
		meta = append(meta, Styles.MutedText.Render(fmt.Sprintf("🌐 %s", d.job.Source)))
	}

	footer := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, actions...),
		"  ",
		lipgloss.JoinHorizontal(lipgloss.Top, meta...),
	)

	return Styles.StatusBar.Width(d.width).Render(footer)
}

func (d *JobDetail) renderEmpty() string {
	return lipgloss.Place(d.width, d.height,
		lipgloss.Center, lipgloss.Center,
		Styles.MutedText.Render("No job selected"),
	)
}

func (d *JobDetail) updateContent() {
	if d.job == nil {
		d.viewport.SetContent("")
		return
	}

	content := d.formatJobContent()
	d.viewport.SetContent(content)
}

func (d *JobDetail) formatJobContent() string {
	var sections []string

	// Description section
	if d.job.Description != "" {
		descTitle := Styles.ListTitle.Render("Description")
		descContent := Styles.Content.Render(d.job.Description)
		sections = append(sections, lipgloss.JoinVertical(lipgloss.Left, descTitle, descContent))
	}

	// Requirements/Traps section
	if len(d.job.Traps) > 0 {
		trapsTitle := Styles.ErrorText.Render("⚠️  Application Requirements/Traps")
		var trapsList []string
		for _, trap := range d.job.Traps {
			trapsList = append(trapsList, fmt.Sprintf("• %s", trap))
		}
		trapsContent := Styles.Content.Render(strings.Join(trapsList, "\n"))
		sections = append(sections, lipgloss.JoinVertical(lipgloss.Left, trapsTitle, trapsContent))
	}

	// Technical details
	details := []string{}
	if d.job.PostedDate.String() != "" {
		details = append(details, fmt.Sprintf("📅 Posted: %s", d.job.PostedDate.Format("Jan 2, 2006")))
	}
	if d.job.URL != "" {
		details = append(details, fmt.Sprintf("🔗 URL: %s", d.job.URL))
	}
	if d.job.JobType != "" {
		details = append(details, fmt.Sprintf("💼 Type: %s", d.job.JobType))
	}

	if len(details) > 0 {
		detailsTitle := Styles.ListTitle.Render("Details")
		detailsContent := Styles.Content.Render(strings.Join(details, "\n"))
		sections = append(sections, lipgloss.JoinVertical(lipgloss.Left, detailsTitle, detailsContent))
	}

	// Application info
	if d.job.Email != "" {
		appTitle := Styles.ListTitle.Render("Application")
		appContent := Styles.SuccessText.Render(fmt.Sprintf("📧 Contact: %s", d.job.Email))
		sections = append(sections, lipgloss.JoinVertical(lipgloss.Left, appTitle, appContent))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
