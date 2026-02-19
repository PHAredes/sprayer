package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"sprayer/internal/profile"
)

// FilterView represents a CHARM-style filter configuration view
type FilterView struct {
	profile      profile.Profile
	width        int
	height       int
	inputs       []textinput.Model
	focusedInput int
}

// NewFilterView creates a new CHARM-style filter view
func NewFilterView(profile profile.Profile) *FilterView {
	inputs := make([]textinput.Model, 0)

	// Create input fields for various filters
	keywordsInput := textinput.New()
	keywordsInput.Placeholder = "golang, rust, typescript"
	keywordsInput.SetValue(strings.Join(profile.Keywords, ", "))
	keywordsInput.Prompt = "Keywords: "
	inputs = append(inputs, keywordsInput)

	excludeInput := textinput.New()
	excludeInput.Placeholder = "senior, lead, manager"
	excludeInput.SetValue(strings.Join(profile.ExcludeKeywords, ", "))
	excludeInput.Prompt = "Exclude: "
	inputs = append(inputs, excludeInput)

	locationsInput := textinput.New()
	locationsInput.Placeholder = "remote, san francisco, new york"
	locationsInput.SetValue(strings.Join(profile.Locations, ", "))
	locationsInput.Prompt = "Locations: "
	inputs = append(inputs, locationsInput)

	companiesInput := textinput.New()
	companiesInput.Placeholder = "google, microsoft, startup"
	companiesInput.SetValue(strings.Join(profile.PreferredCompanies, ", "))
	companiesInput.Prompt = "Companies: "
	inputs = append(inputs, companiesInput)

	minScoreInput := textinput.New()
	minScoreInput.Placeholder = "0"
	minScoreInput.SetValue(fmt.Sprintf("%d", profile.MinScore))
	minScoreInput.Prompt = "Min Score: "
	inputs = append(inputs, minScoreInput)

	// Focus the first input
	if len(inputs) > 0 {
		inputs[0].Focus()
	}

	return &FilterView{
		profile:      profile,
		inputs:       inputs,
		focusedInput: 0,
	}
}

// SetSize updates the component size
func (f *FilterView) SetSize(width, height int) {
	f.width = width
	f.height = height

	// Update input widths
	for i := range f.inputs {
		f.inputs[i].Width = width - 20
	}
}

// Update handles messages
func (f *FilterView) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			f.nextInput()
		case "shift+tab", "up":
			f.prevInput()
		case "enter":
			// Apply filters
			return func() tea.Msg {
				return FiltersAppliedMsg{Profile: f.collectInput()}
			}
		}
	}

	// Update focused input
	if f.focusedInput >= 0 && f.focusedInput < len(f.inputs) {
		var cmd tea.Cmd
		f.inputs[f.focusedInput], cmd = f.inputs[f.focusedInput].Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

// View renders the component
func (f *FilterView) View(width, height int) string {
	f.SetSize(width, height)

	header := f.renderHeader()
	form := f.renderForm()
	footer := f.renderFooter()

	contentHeight := height - 6 // Account for header and footer
	formContent := lipgloss.NewStyle().
		Height(contentHeight).
		Width(width - 4).
		Render(form)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		formContent,
		footer,
	)
}

func (f *FilterView) renderHeader() string {
	title := Styles.Title.Render("Filter Configuration")
	help := Styles.MutedText.Render("Tab: next field • Enter: apply • Esc: cancel")

	return lipgloss.JoinHorizontal(lipgloss.Top,
		Styles.Header.Width(f.width/2).Render(title),
		Styles.Header.Width(f.width/2).Align(lipgloss.Right).Render(help),
	)
}

func (f *FilterView) renderForm() string {
	var rows []string

	for i, input := range f.inputs {
		row := f.renderInputRow(input, i == f.focusedInput)
		rows = append(rows, row)
		rows = append(rows, "") // Spacing between inputs
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (f *FilterView) renderInputRow(input textinput.Model, focused bool) string {
	prompt := Styles.InputPrompt.Render(input.Prompt)

	if focused {
		return lipgloss.JoinHorizontal(lipgloss.Top,
			prompt,
			Styles.InputFocused.Width(f.width-20).Render(input.View()),
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		prompt,
		Styles.Input.Width(f.width-20).Render(input.View()),
	)
}

func (f *FilterView) renderFooter() string {
	// Show current profile info
	profileInfo := fmt.Sprintf("Profile: %s", f.profile.Name)

	// Quick filter toggles
	toggles := []string{}
	if f.profile.ExcludeTraps {
		toggles = append(toggles, Styles.SuccessText.Render("✓ No traps"))
	}
	if f.profile.MustHaveEmail {
		toggles = append(toggles, Styles.SuccessText.Render("✓ Must have email"))
	}
	if f.profile.PreferRemote {
		toggles = append(toggles, Styles.SuccessText.Render("✓ Remote preferred"))
	}

	left := Styles.StatusText.Render(profileInfo)
	right := Styles.StatusText.Render(strings.Join(toggles, "  "))

	return lipgloss.JoinHorizontal(lipgloss.Top,
		Styles.StatusBar.Width(f.width/2).Render(left),
		Styles.StatusBar.Width(f.width/2).Align(lipgloss.Right).Render(right),
	)
}

func (f *FilterView) nextInput() {
	f.inputs[f.focusedInput].Blur()
	f.focusedInput = (f.focusedInput + 1) % len(f.inputs)
	f.inputs[f.focusedInput].Focus()
}

func (f *FilterView) prevInput() {
	f.inputs[f.focusedInput].Blur()
	f.focusedInput--
	if f.focusedInput < 0 {
		f.focusedInput = len(f.inputs) - 1
	}
	f.inputs[f.focusedInput].Focus()
}

func (f *FilterView) collectInput() profile.Profile {
	// Collect input values and update profile
	updatedProfile := f.profile

	if len(f.inputs) > 0 {
		keywords := strings.Split(f.inputs[0].Value(), ",")
		for i := range keywords {
			keywords[i] = strings.TrimSpace(keywords[i])
		}
		updatedProfile.Keywords = keywords
	}

	if len(f.inputs) > 1 {
		excludeKeywords := strings.Split(f.inputs[1].Value(), ",")
		for i := range excludeKeywords {
			excludeKeywords[i] = strings.TrimSpace(excludeKeywords[i])
		}
		updatedProfile.ExcludeKeywords = excludeKeywords
	}

	if len(f.inputs) > 2 {
		locations := strings.Split(f.inputs[2].Value(), ",")
		for i := range locations {
			locations[i] = strings.TrimSpace(locations[i])
		}
		updatedProfile.Locations = locations
	}

	if len(f.inputs) > 3 {
		companies := strings.Split(f.inputs[3].Value(), ",")
		for i := range companies {
			companies[i] = strings.TrimSpace(companies[i])
		}
		updatedProfile.PreferredCompanies = companies
	}

	if len(f.inputs) > 4 {
		var minScore int
		fmt.Sscanf(f.inputs[4].Value(), "%d", &minScore)
		updatedProfile.MinScore = minScore
	}

	return updatedProfile
}

// FiltersAppliedMsg is handled by the main app (defined in app.go)
