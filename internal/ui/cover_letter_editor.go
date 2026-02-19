package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/internal/job"
)

type EditorState int

const (
	EditorStateIdle EditorState = iota
	EditorStateGenerating
	EditorStateEditing
	EditorStateSending
)

type CoverLetterEditor struct {
	textarea     textarea.Model
	job          *job.Job
	coverLetter  string
	width        int
	height       int
	focusElement focusElement
	state        EditorState
	err          error
	charCount    int
	wordCount    int
}

type focusElement int

const (
	focusTextarea focusElement = iota
)

func NewCoverLetterEditor() *CoverLetterEditor {
	ta := textarea.New()
	ta.Placeholder = "Cover letter will be generated here..."
	ta.SetHeight(15)
	ta.ShowLineNumbers = false
	ta.CharLimit = 5000

	return &CoverLetterEditor{
		textarea:     ta,
		focusElement: focusTextarea,
		state:        EditorStateIdle,
	}
}

func (e *CoverLetterEditor) SetJob(j *job.Job) {
	e.job = j
	e.err = nil
}

func (e *CoverLetterEditor) SetContent(content string) {
	e.coverLetter = content
	e.textarea.SetValue(content)
	e.state = EditorStateEditing
	e.updateCounts()
}

func (e *CoverLetterEditor) Content() string {
	return e.textarea.Value()
}

func (e *CoverLetterEditor) SetSize(width, height int) {
	e.width = width
	e.height = height
	e.textarea.SetWidth(width - 6)
	e.textarea.SetHeight(height - 12)
}

func (e *CoverLetterEditor) SetGenerating() {
	e.state = EditorStateGenerating
	e.textarea.SetValue("Generating cover letter...")
}

func (e *CoverLetterEditor) SetError(err error) {
	e.err = err
	e.state = EditorStateIdle
	e.textarea.SetValue("")
}

func (e *CoverLetterEditor) State() EditorState {
	return e.state
}

func (e *CoverLetterEditor) Init() tea.Cmd {
	return textarea.Blink
}

func (e *CoverLetterEditor) updateCounts() {
	content := e.textarea.Value()
	e.charCount = len(content)
	e.wordCount = countWords(content)
}

func countWords(s string) int {
	count := 0
	inWord := false
	for _, r := range s {
		if r == ' ' || r == '\n' || r == '\t' {
			inWord = false
		} else if !inWord {
			inWord = true
			count++
		}
	}
	return count
}

func (e *CoverLetterEditor) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if e.state == EditorStateGenerating {
			return nil
		}

		switch {
		case key.Matches(msg, CoverLetterKeys.Back):
			return nil
		case key.Matches(msg, CoverLetterKeys.Send):
			if e.state == EditorStateEditing && e.textarea.Value() != "" {
				e.state = EditorStateSending
				return func() tea.Msg {
					return CoverLetterSendMsg{
						JobID:   e.job.ID,
						Content: e.textarea.Value(),
					}
				}
			}
			return nil
		case key.Matches(msg, CoverLetterKeys.Regenerate):
			if e.job != nil {
				return func() tea.Msg {
					return CoverLetterRegenerateMsg{JobID: e.job.ID}
				}
			}
			return nil
		case key.Matches(msg, CoverLetterKeys.Save):
			if e.state == EditorStateEditing {
				return func() tea.Msg {
					return CoverLetterSavedMsg{
						JobID:   e.job.ID,
						Content: e.textarea.Value(),
					}
				}
			}
			return nil
		}
	}

	var cmd tea.Cmd
	e.textarea, cmd = e.textarea.Update(msg)
	cmds = append(cmds, cmd)
	e.updateCounts()

	return tea.Batch(cmds...)
}

func (e *CoverLetterEditor) View() string {
	if e.job == nil {
		return "No job selected"
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(Colors.Accent).
		Padding(0, 1)

	var statusText string
	switch e.state {
	case EditorStateGenerating:
		statusText = lipgloss.NewStyle().Foreground(Colors.Warning).Render("⏳ Generating...")
	case EditorStateEditing:
		statusText = lipgloss.NewStyle().Foreground(Colors.Success).Render("✏️  Editing")
	case EditorStateSending:
		statusText = lipgloss.NewStyle().Foreground(Colors.Primary).Render("📤 Sending...")
	default:
		statusText = lipgloss.NewStyle().Foreground(Colors.Muted).Render("Ready")
	}

	header := headerStyle.Render(fmt.Sprintf("Cover Letter for %s at %s", e.job.Title, e.job.Company))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Surface).
		Padding(1)

	if e.state == EditorStateGenerating {
		boxStyle = boxStyle.BorderForeground(Colors.Warning)
	} else if e.state == EditorStateEditing {
		boxStyle = boxStyle.BorderForeground(Colors.Success)
	}

	textareaBox := boxStyle.Width(e.width - 4).Render(e.textarea.View())

	statsStyle := lipgloss.NewStyle().
		Foreground(Colors.Muted).
		Padding(0, 1)

	stats := statsStyle.Render(fmt.Sprintf("%d chars • %d words", e.charCount, e.wordCount))

	helpStyle := lipgloss.NewStyle().
		Foreground(Colors.Muted).
		Padding(1, 0)

	help := helpStyle.Render("Ctrl+Enter: Send • Ctrl+S: Save • Ctrl+R: Regenerate • Esc: Back")

	var errorView string
	if e.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(Colors.Error).
			Padding(1, 0)
		errorView = errorStyle.Render(fmt.Sprintf("⚠️  Error: %v", e.err))
	}

	sections := []string{header, statusText, textareaBox, stats, help}
	if errorView != "" {
		sections = append(sections, errorView)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

var CoverLetterKeys = coverLetterKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Send: key.NewBinding(
		key.WithKeys("ctrl+enter"),
		key.WithHelp("ctrl+enter", "send application"),
	),
	Regenerate: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "regenerate"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	),
}

type coverLetterKeyMap struct {
	Back       key.Binding
	Send       key.Binding
	Regenerate key.Binding
	Save       key.Binding
}

func (k coverLetterKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Send, k.Regenerate, k.Save}
}

func (k coverLetterKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Send},
		{k.Regenerate, k.Save},
	}
}

type CoverLetterGeneratedMsg struct {
	Content string
	JobID   string
}

type CoverLetterRegenerateMsg struct {
	JobID string
}

type CoverLetterSavedMsg struct {
	Content string
	JobID   string
}

type CoverLetterSendMsg struct {
	Content string
	JobID   string
}

type CoverLetterError struct {
	Error error
	JobID string
}
