package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"sprayer/internal/apply"
	"sprayer/internal/job"
	"sprayer/internal/profile"
)

type CVChoice int

const (
	CVChoiceNone CVChoice = iota
	CVChoiceOriginal
	CVChoiceCustom
)

type CVPreviewView struct {
	job           *job.Job
	profile       *profile.Profile
	cvGenerator   *apply.CVGenerator
	viewport      viewport.Model
	width         int
	height        int
	choice        CVChoice
	customCV      string
	originalCV    string
	generating    bool
	generateError string
}

func NewCVPreviewView() *CVPreviewView {
	vp := viewport.New(80, 20)
	vp.Style = Styles.Detail

	return &CVPreviewView{
		viewport: vp,
		choice:   CVChoiceNone,
	}
}

func (v *CVPreviewView) SetJob(j *job.Job) {
	v.job = j
	v.choice = CVChoiceNone
	v.customCV = ""
	v.generateError = ""
	v.generating = false
	v.updateContent()
}

func (v *CVPreviewView) SetProfile(p *profile.Profile) {
	v.profile = p
	if p != nil && p.CVData != nil {
		v.originalCV = v.formatOriginalCV(p.CVData)
	}
}

func (v *CVPreviewView) SetCVGenerator(gen *apply.CVGenerator) {
	v.cvGenerator = gen
}

func (v *CVPreviewView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.viewport.Width = width - 4
	v.viewport.Height = height - 10
}

func (v *CVPreviewView) Init() tea.Cmd {
	return nil
}

func (v *CVPreviewView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			v.choice = CVChoiceOriginal
			v.updateContent()
		case "2":
			if v.customCV == "" && !v.generating && v.cvGenerator != nil {
				return v.generateCustomCV()
			}
			v.choice = CVChoiceCustom
			v.updateContent()
		case "up", "k":
			v.viewport.LineUp(1)
		case "down", "j":
			v.viewport.LineDown(1)
		}
	case CVGeneratedMsg:
		v.generating = false
		if msg.Error != nil {
			v.generateError = msg.Error.Error()
		} else {
			v.customCV = msg.Content
			v.choice = CVChoiceCustom
			v.updateContent()
		}
	}

	v.viewport, cmd = v.viewport.Update(msg)
	return cmd
}

func (v *CVPreviewView) generateCustomCV() tea.Cmd {
	return func() tea.Msg {
		if v.cvGenerator == nil || v.job == nil || v.profile == nil {
			return CVGeneratedMsg{Error: fmt.Errorf("missing required data")}
		}

		content, err := v.cvGenerator.GenerateCustomCV(v.job, v.profile)
		if err != nil {
			return CVGeneratedMsg{Error: err}
		}
		return CVGeneratedMsg{Content: content}
	}
}

func (v *CVPreviewView) View() string {
	if v.job == nil {
		return Styles.MutedText.Render("No job selected")
	}

	var sections []string

	header := v.renderHeader()
	sections = append(sections, header)

	if v.generating {
		loading := lipgloss.NewStyle().
			Foreground(Colors.Accent).
			Render("Generating custom CV...")
		sections = append(sections, loading)
	} else if v.generateError != "" {
		errText := Styles.ErrorText.Render(fmt.Sprintf("Error: %s", v.generateError))
		sections = append(sections, errText)
	}

	sections = append(sections, v.viewport.View())

	footer := v.renderFooter()
	sections = append(sections, footer)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (v *CVPreviewView) renderHeader() string {
	title := Styles.DetailTitle.Render(fmt.Sprintf("CV Preview - %s @ %s", v.job.Title, v.job.Company))

	var choiceText string
	switch v.choice {
	case CVChoiceOriginal:
		choiceText = Styles.SuccessText.Render("Using Original CV")
	case CVChoiceCustom:
		choiceText = Styles.SuccessText.Render("Using Custom CV")
	default:
		choiceText = Styles.MutedText.Render("Select CV type")
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, choiceText)
}

func (v *CVPreviewView) renderFooter() string {
	var options []string
	options = append(options, Styles.Text.Render("[1] Original CV"))
	options = append(options, Styles.Text.Render("[2] Generate Custom CV"))

	if v.choice != CVChoiceNone {
		options = append(options, Styles.SuccessText.Render("[Enter] Apply with selected CV"))
	}
	options = append(options, Styles.MutedText.Render("[Esc] Cancel"))

	footer := lipgloss.JoinHorizontal(lipgloss.Top, options...)
	return Styles.StatusBar.Width(v.width).Render(footer)
}

func (v *CVPreviewView) updateContent() {
	var content string

	switch v.choice {
	case CVChoiceOriginal:
		content = v.originalCV
	case CVChoiceCustom:
		content = v.customCV
	default:
		content = v.renderChoicePrompt()
	}

	v.viewport.SetContent(content)
}

func (v *CVPreviewView) renderChoicePrompt() string {
	var b strings.Builder
	b.WriteString(Styles.Title.Render("Choose CV Type") + "\n\n")

	b.WriteString("[1] Use Original CV\n")
	b.WriteString("    Your standard CV from profile\n\n")

	b.WriteString("[2] Generate Custom CV\n")
	if v.cvGenerator != nil && v.cvGenerator.Available() {
		b.WriteString("    AI-tailored CV for this specific job\n")
	} else {
		b.WriteString(Styles.MutedText.Render("    (LLM not configured - set SPRAYER_LLM_KEY)\n"))
	}

	return b.String()
}

func (v *CVPreviewView) formatOriginalCV(cv *profile.CVData) string {
	var b strings.Builder

	if cv.Name != "" {
		b.WriteString(Styles.Title.Render(cv.Name) + "\n")
	}
	if cv.Title != "" {
		b.WriteString(Styles.Subtitle.Render(cv.Title) + "\n")
	}
	b.WriteString("\n")

	if cv.Email != "" || cv.Phone != "" {
		b.WriteString("Contact:\n")
		if cv.Email != "" {
			b.WriteString(fmt.Sprintf("  Email: %s\n", cv.Email))
		}
		if cv.Phone != "" {
			b.WriteString(fmt.Sprintf("  Phone: %s\n", cv.Phone))
		}
		b.WriteString("\n")
	}

	if cv.Summary != "" {
		b.WriteString("Summary:\n")
		b.WriteString(fmt.Sprintf("%s\n\n", cv.Summary))
	}

	if len(cv.Technologies) > 0 {
		b.WriteString("Technologies:\n")
		b.WriteString(fmt.Sprintf("  %s\n\n", strings.Join(cv.Technologies, ", ")))
	}

	if len(cv.Skills) > 0 {
		b.WriteString("Skills:\n")
		b.WriteString(fmt.Sprintf("  %s\n\n", strings.Join(cv.Skills, ", ")))
	}

	if len(cv.Experience) > 0 {
		b.WriteString("Experience:\n")
		for _, exp := range cv.Experience {
			b.WriteString(fmt.Sprintf("\n  %s at %s", exp.Title, exp.Company))
			if exp.Duration != "" {
				b.WriteString(fmt.Sprintf(" (%s)", exp.Duration))
			}
			b.WriteString("\n")
			if exp.Description != "" {
				b.WriteString(fmt.Sprintf("    %s\n", exp.Description))
			}
		}
		b.WriteString("\n")
	}

	if len(cv.Education) > 0 {
		b.WriteString("Education:\n")
		for _, edu := range cv.Education {
			b.WriteString(fmt.Sprintf("\n  %s in %s - %s", edu.Degree, edu.Field, edu.Institution))
			if edu.Year != "" {
				b.WriteString(fmt.Sprintf(" (%s)", edu.Year))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (v *CVPreviewView) SelectedCV() string {
	switch v.choice {
	case CVChoiceOriginal:
		return v.originalCV
	case CVChoiceCustom:
		return v.customCV
	default:
		return ""
	}
}

func (v *CVPreviewView) HasSelection() bool {
	return v.choice != CVChoiceNone && v.SelectedCV() != ""
}

type CVGeneratedMsg struct {
	Content string
	Error   error
}
