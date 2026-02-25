package cvwizard

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/theme"
)

// WizardStep represents the current step in the CV wizard.
type WizardStep int

const (
	StepEntry WizardStep = iota
	StepInfo
	StepSummary
	StepExperience
	StepSkills
	StepReview
)

// ExperienceEntry holds one work experience record.
type ExperienceEntry struct {
	Role       string
	Company    string
	Period     string
	Stack      []string
	Highlights []string
}

// DoneMsg signals the wizard completed successfully.
type DoneMsg struct {
	ProfileName string
}

// CancelMsg signals the wizard was cancelled.
type CancelMsg struct{}

// Model holds all CV wizard state across steps.
type Model struct {
	Step       WizardStep
	ImportMode bool   // false = manual, true = import
	FilePath   string // only used in import mode

	// Step 1: Personal Info
	Name     string
	JobTitle string
	Email    string
	GitHub   string
	LinkedIn string
	Location string

	// Step 2: Summary
	Summary string

	// Step 3: Experience
	Experiences []ExperienceEntry
	CurrentExp  ExperienceEntry
	CurrentTag  string // tag being typed for stack

	// Step 4: Skills
	Languages  []string
	Frameworks []string
	Databases  []string
	Seniority  []string
	SkillInput string // current tag being typed

	// Step 5: Review
	ProfileName string

	// UI state
	FocusIndex int
	Width      int
	Height     int
}

// New creates a fresh CV wizard.
func New(w, h int) Model {
	return Model{
		Step:        StepEntry,
		ImportMode:  false,
		Experiences: []ExperienceEntry{},
		Languages:   []string{},
		Frameworks:  []string{},
		Databases:   []string{},
		Seniority:   []string{},
		ProfileName: "Default",
		FocusIndex:  0,
		Width:       w,
		Height:      h,
	}
}

// Update dispatches to the current step's update handler.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.Step {
	case StepEntry:
		return m.updateEntry(msg)
	case StepInfo:
		return m.updateInfo(msg)
	case StepSummary:
		return m.updateSummary(msg)
	case StepExperience:
		return m.updateExperience(msg)
	case StepSkills:
		return m.updateSkills(msg)
	case StepReview:
		return m.updateReview(msg)
	}
	return m, nil
}

// View dispatches to the current step's view.
func (m Model) View() string {
	switch m.Step {
	case StepEntry:
		return m.viewEntry()
	case StepInfo:
		return m.viewInfo()
	case StepSummary:
		return m.viewSummary()
	case StepExperience:
		return m.viewExperience()
	case StepSkills:
		return m.viewSkills()
	case StepReview:
		return m.viewReview()
	}
	return ""
}

// ModalTopBar renders the breadcrumb top bar for non-entry steps.
func (m Model) ModalTopBar() string {
	if m.Step == StepEntry {
		return ""
	}

	crumbs := m.breadcrumbs()
	hint := m.hintText()

	crumbStr := ""
	for i, c := range crumbs {
		if i > 0 {
			crumbStr += theme.CVBreadcrumbSepStyle.Background(theme.Surface2).Render(" > ")
		}
		crumbStr += c
	}

	hintStr := theme.ModalHintStyle.Render(hint)

	crumbW := lipgloss.Width(crumbStr)
	hintW := lipgloss.Width(hintStr)
	gap := m.Width - crumbW - hintW - 4
	if gap < 1 {
		gap = 1
	}
	pad := lipgloss.NewStyle().Background(theme.Surface2).Width(gap).Render("")

	row := crumbStr + pad + hintStr
	return theme.ModalTopBarStyle.Width(m.Width).PaddingLeft(2).PaddingRight(2).Render(row)
}

func (m Model) breadcrumbs() []string {
	type crumb struct {
		Label string
		Step  WizardStep
	}

	allCrumbs := []crumb{
		{"1. Personal Info", StepInfo},
		{"2. Summary", StepSummary},
		{"3. Experience", StepExperience},
		{"4. Skills", StepSkills},
		{"5. Review", StepReview},
	}

	// Show a window of 3 crumbs around current step
	currentIdx := 0
	for i, c := range allCrumbs {
		if c.Step == m.Step {
			currentIdx = i
			break
		}
	}

	start := currentIdx - 1
	if start < 0 {
		start = 0
	}
	end := start + 3
	if end > len(allCrumbs) {
		end = len(allCrumbs)
		start = end - 3
		if start < 0 {
			start = 0
		}
	}

	var result []string
	for _, c := range allCrumbs[start:end] {
		var style lipgloss.Style
		if c.Step == m.Step {
			style = theme.CVBreadcrumbCurrentStyle.
				PaddingLeft(1).PaddingRight(1)
		} else if c.Step < m.Step {
			style = theme.CVBreadcrumbDoneStyle.Background(theme.Surface2)
		} else {
			style = theme.CVBreadcrumbStyle.Background(theme.Surface2)
		}
		result = append(result, style.Render(c.Label))
	}

	return result
}

func (m Model) hintText() string {
	switch m.Step {
	case StepInfo:
		return "tab next field  |  enter continue"
	case StepSummary:
		return "ctrl+enter continue  |  esc back"
	case StepExperience:
		return "ctrl+n add another  |  ctrl+enter next  |  esc back"
	case StepSkills:
		return "enter add tag  |  ctrl+enter continue"
	case StepReview:
		return "enter save  |  esc back"
	}
	return ""
}

// StepDots renders step progress dots for the status bar.
func (m Model) StepDots() string {
	total := 5
	current := int(m.Step)
	if current < 1 {
		current = 1
	}

	var dots []string
	for i := 1; i <= total; i++ {
		if i <= current {
			dots = append(dots, theme.CVStepDotActiveStyle.Render("●"))
		} else {
			dots = append(dots, theme.CVStepDotInactiveStyle.Render("●"))
		}
	}
	return strings.Join(dots, " ")
}

// collectKeywords gathers all skills for the profile.
func (m Model) collectKeywords() []string {
	var kw []string
	kw = append(kw, m.Languages...)
	kw = append(kw, m.Frameworks...)
	kw = append(kw, m.Databases...)
	return kw
}
