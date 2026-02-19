package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"sprayer/internal/profile"
)

// ProfileView represents a CHARM-style profile management view
type ProfileView struct {
	profiles      []profile.Profile
	activeProfile profile.Profile
	cursor        int
	width         int
	height        int
}

// NewProfileView creates a new CHARM-style profile view
func NewProfileView(profiles []profile.Profile, activeProfile profile.Profile) *ProfileView {
	cursor := 0
	for i, p := range profiles {
		if p.ID == activeProfile.ID {
			cursor = i
			break
		}
	}

	return &ProfileView{
		profiles:      profiles,
		activeProfile: activeProfile,
		cursor:        cursor,
	}
}

// SelectedProfile returns the currently selected profile
func (p *ProfileView) SelectedProfile() profile.Profile {
	if p.cursor >= 0 && p.cursor < len(p.profiles) {
		return p.profiles[p.cursor]
	}
	return p.activeProfile
}

// MoveCursor moves the cursor up or down
func (p *ProfileView) MoveCursor(up bool) {
	if up && p.cursor > 0 {
		p.cursor--
	} else if !up && p.cursor < len(p.profiles)-1 {
		p.cursor++
	}
}

// SetSize updates the component size
func (p *ProfileView) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// Update handles messages
func (p *ProfileView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Up):
			p.MoveCursor(true)
		case key.Matches(msg, Keys.Down):
			p.MoveCursor(false)
		}
	}
	return nil
}

// View renders the component
func (p *ProfileView) View(width, height int) string {
	p.SetSize(width, height)

	header := p.renderHeader()
	list := p.renderProfileList()
	details := p.renderProfileDetails()

	// Split view: list on left, details on right
	leftWidth := width / 3
	rightWidth := width - leftWidth - 1

	leftPane := lipgloss.NewStyle().
		Width(leftWidth).
		Height(height - 3).
		Render(list)

	rightPane := lipgloss.NewStyle().
		Width(rightWidth).
		Height(height - 3).
		Render(details)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	return lipgloss.JoinVertical(lipgloss.Left, header, content)
}

func (p *ProfileView) renderHeader() string {
	title := Styles.Title.Render("Profiles")
	help := Styles.MutedText.Render("↑/↓: select • Enter: apply • Esc: back")

	return lipgloss.JoinHorizontal(lipgloss.Top,
		Styles.Header.Width(p.width/2).Render(title),
		Styles.Header.Width(p.width/2).Align(lipgloss.Right).Render(help),
	)
}

func (p *ProfileView) renderProfileList() string {
	var rows []string

	for i, profile := range p.profiles {
		row := p.renderProfileRow(profile, i == p.cursor)
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return Styles.MutedText.Render("No profiles available")
	}

	return strings.Join(rows, "\n")
}

func (p *ProfileView) renderProfileRow(profile profile.Profile, selected bool) string {
	name := profile.Name
	if profile.ID == p.activeProfile.ID {
		name = fmt.Sprintf("%s (active)", name)
	}

	// Filter summary
	filterSummary := profile.GetFilterSummary()
	if len(filterSummary) > 30 {
		filterSummary = filterSummary[:27] + "..."
	}

	row := fmt.Sprintf("%-20s %s", name, filterSummary)

	if selected {
		return Styles.SelectedItem.Width(p.width / 3).Render(row)
	}

	return Styles.ListItem.Width(p.width / 3).Render(row)
}

func (p *ProfileView) renderProfileDetails() string {
	if p.cursor >= len(p.profiles) {
		return Styles.MutedText.Render("No profile selected")
	}

	profile := p.profiles[p.cursor]

	// Build details sections
	var sections []string

	// Basic info
	basicInfo := []string{
		Styles.DetailTitle.Render("Profile Information"),
		fmt.Sprintf("Name: %s", profile.Name),
		fmt.Sprintf("Keywords: %s", strings.Join(profile.Keywords, ", ")),
	}

	if profile.PreferRemote {
		basicInfo = append(basicInfo, "Remote: Preferred")
	}

	if len(profile.Locations) > 0 {
		basicInfo = append(basicInfo, fmt.Sprintf("Locations: %s", strings.Join(profile.Locations, ", ")))
	}

	sections = append(sections, strings.Join(basicInfo, "\n"))

	// Filtering preferences
	filterInfo := []string{
		Styles.DetailTitle.Render("Filtering Preferences"),
		fmt.Sprintf("Score Range: %d-%d", profile.MinScore, profile.MaxScore),
	}

	if profile.ExcludeTraps {
		filterInfo = append(filterInfo, "Exclude Traps: Yes")
	}

	if profile.MustHaveEmail {
		filterInfo = append(filterInfo, "Must Have Email: Yes")
	}

	if len(profile.SeniorityLevels) > 0 {
		filterInfo = append(filterInfo, fmt.Sprintf("Seniority: %s", strings.Join(profile.SeniorityLevels, ", ")))
	}

	if len(profile.PreferredTech) > 0 {
		filterInfo = append(filterInfo, fmt.Sprintf("Preferred Tech: %s", strings.Join(profile.PreferredTech, ", ")))
	}

	if len(profile.AvoidTech) > 0 {
		filterInfo = append(filterInfo, fmt.Sprintf("Avoid Tech: %s", strings.Join(profile.AvoidTech, ", ")))
	}

	sections = append(sections, strings.Join(filterInfo, "\n"))

	// Scoring weights
	weightsInfo := []string{
		Styles.DetailTitle.Render("Scoring Weights"),
		fmt.Sprintf("Tech Match: %d%%", profile.ScoringWeights.TechMatch),
		fmt.Sprintf("Seniority Match: %d%%", profile.ScoringWeights.SeniorityMatch),
		fmt.Sprintf("Location Match: %d%%", profile.ScoringWeights.LocationMatch),
		fmt.Sprintf("Company Match: %d%%", profile.ScoringWeights.CompanyMatch),
		fmt.Sprintf("Salary Match: %d%%", profile.ScoringWeights.SalaryMatch),
		fmt.Sprintf("Remote Match: %d%%", profile.ScoringWeights.RemoteMatch),
	}

	sections = append(sections, strings.Join(weightsInfo, "\n"))

	return Styles.Detail.Width(p.width * 2 / 3).Render(
		lipgloss.JoinVertical(lipgloss.Left, sections...),
	)
}
