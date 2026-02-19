package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"sprayer/internal/apply"
	"sprayer/internal/job"
	"sprayer/internal/profile"
)

type EmailComposerField int

const (
	FieldTo EmailComposerField = iota
	FieldSubject
	FieldBody
)

type EmailComposer struct {
	draft        *apply.EmailDraft
	job          *job.Job
	profile      profile.Profile
	toInput      textinput.Model
	subjectInput textinput.Model
	body         textarea.Model
	focused      int
	width        int
	height       int
	err          error
	result       EmailComposerResult
	sendCmd      func() error
	cancelCmd    func() error
}

type EmailComposerResult struct {
	Action    EmailAction
	Draft     *apply.EmailDraft
	Cancelled bool
}

type EmailAction int

const (
	ActionSend EmailAction = iota
	ActionSaveDraft
	ActionCancel
)

func NewEmailComposer(j *job.Job, p profile.Profile, emailContent string) *EmailComposer {
	toInput := textinput.New()
	toInput.Placeholder = "recipient@example.com"
	toInput.CharLimit = 256
	if j != nil && j.Email != "" {
		toInput.SetValue(j.Email)
	}

	subjectInput := textinput.New()
	subjectInput.Placeholder = "Email subject..."
	subjectInput.CharLimit = 200

	body := textarea.New()
	body.Placeholder = "Write your email body here..."
	body.SetHeight(10)
	body.ShowLineNumbers = false
	body.CharLimit = 50000
	body.SetValue(emailContent)

	return &EmailComposer{
		job:          j,
		profile:      p,
		toInput:      toInput,
		subjectInput: subjectInput,
		body:         body,
		focused:      0,
	}
}

func (e *EmailComposer) SetDraft(draft *apply.EmailDraft) {
	e.draft = draft
	if draft != nil {
		e.toInput.SetValue(draft.To)
		e.subjectInput.SetValue(draft.Subject)
		e.body.SetValue(draft.Body)
	}
}

func (e *EmailComposer) SetSize(width, height int) {
	e.width = width
	e.height = height
	e.body.SetWidth(width - 8)
	e.body.SetHeight(height - 16)
}

func (e *EmailComposer) GetDraft() *apply.EmailDraft {
	if e.draft == nil {
		e.draft = &apply.EmailDraft{
			ID:        fmt.Sprintf("%s.sprayer", e.job.ID),
			From:      e.profile.ContactEmail,
			CreatedAt: e.job.PostedDate,
		}
		if e.job != nil && len(e.draft.AttachmentPaths) == 0 {
			cvPDF := findPDF(e.profile.CVPath)
			if cvPDF != "" {
				e.draft.AttachmentPaths = []string{cvPDF}
			}
		}
	}
	e.draft.To = e.toInput.Value()
	e.draft.Subject = e.subjectInput.Value()
	e.draft.Body = e.body.Value()
	return e.draft
}

func (e *EmailComposer) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, textarea.Blink)
}

func (e *EmailComposer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, EmailComposerKeys.Quit):
			e.result = EmailComposerResult{Action: ActionCancel, Cancelled: true}
			return e, tea.Quit

		case key.Matches(msg, EmailComposerKeys.Cancel):
			e.result = EmailComposerResult{Action: ActionCancel, Cancelled: true}
			return e, tea.Quit

		case key.Matches(msg, EmailComposerKeys.Send):
			if e.validate() {
				e.result = EmailComposerResult{Action: ActionSend, Draft: e.GetDraft()}
				return e, tea.Quit
			}

		case key.Matches(msg, EmailComposerKeys.SaveDraft):
			e.result = EmailComposerResult{Action: ActionSaveDraft, Draft: e.GetDraft()}
			return e, tea.Quit

		case key.Matches(msg, EmailComposerKeys.Tab):
			e.focusNext()

		case key.Matches(msg, EmailComposerKeys.Backtab):
			e.focusPrev()

		default:
			switch e.focused {
			case int(FieldTo):
				var cmd tea.Cmd
				e.toInput, cmd = e.toInput.Update(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			case int(FieldSubject):
				var cmd tea.Cmd
				e.subjectInput, cmd = e.subjectInput.Update(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			case int(FieldBody):
				var cmd tea.Cmd
				e.body, cmd = e.body.Update(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		}
	}

	return e, tea.Batch(cmds...)
}

func (e *EmailComposer) focusNext() {
	e.focused = (e.focused + 1) % 3
	e.updateFocus()
}

func (e *EmailComposer) focusPrev() {
	e.focused = (e.focused + 2) % 3
	e.updateFocus()
}

func (e *EmailComposer) updateFocus() {
	switch e.focused {
	case int(FieldTo):
		e.toInput.Focus()
		e.subjectInput.Blur()
		e.body.Blur()
	case int(FieldSubject):
		e.toInput.Blur()
		e.subjectInput.Focus()
		e.body.Blur()
	case int(FieldBody):
		e.toInput.Blur()
		e.subjectInput.Blur()
		e.body.Focus()
	}
}

func (e *EmailComposer) validate() bool {
	e.err = nil
	if strings.TrimSpace(e.toInput.Value()) == "" {
		e.err = fmt.Errorf("recipient (To) is required")
		return false
	}
	if strings.TrimSpace(e.subjectInput.Value()) == "" {
		e.err = fmt.Errorf("subject is required")
		return false
	}
	if strings.TrimSpace(e.body.Value()) == "" {
		e.err = fmt.Errorf("body is required")
		return false
	}
	return true
}

func (e *EmailComposer) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(Colors.Primary).
		Padding(0, 1).
		MarginBottom(1)

	jobInfo := ""
	if e.job != nil {
		jobInfo = fmt.Sprintf(" - %s @ %s", e.job.Title, e.job.Company)
	}
	b.WriteString(titleStyle.Render("Compose Email" + jobInfo))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().
		Foreground(Colors.Subtle).
		Width(10)

	focusedLabelStyle := lipgloss.NewStyle().
		Foreground(Colors.Accent).
		Bold(true).
		Width(10)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Border).
		Padding(0, 1).
		MarginLeft(2)

	focusedInputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Accent).
		Padding(0, 1).
		MarginLeft(2)

	toLabel := labelStyle.Render("To:")
	if e.focused == int(FieldTo) {
		toLabel = focusedLabelStyle.Render("To:")
		e.toInput.Focus()
		toBox := focusedInputStyle.Width(e.width - 16).Render(e.toInput.View())
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, toLabel, toBox))
	} else {
		e.toInput.Blur()
		toBox := inputStyle.Width(e.width - 16).Render(e.toInput.View())
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, toLabel, toBox))
	}
	b.WriteString("\n\n")

	subjectLabel := labelStyle.Render("Subject:")
	if e.focused == int(FieldSubject) {
		subjectLabel = focusedLabelStyle.Render("Subject:")
		e.subjectInput.Focus()
		subjectBox := focusedInputStyle.Width(e.width - 16).Render(e.subjectInput.View())
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, subjectLabel, subjectBox))
	} else {
		e.subjectInput.Blur()
		subjectBox := inputStyle.Width(e.width - 16).Render(e.subjectInput.View())
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, subjectLabel, subjectBox))
	}
	b.WriteString("\n\n")

	bodyLabel := labelStyle.Render("Body:")
	if e.focused == int(FieldBody) {
		bodyLabel = focusedLabelStyle.Render("Body:")
	}
	b.WriteString(bodyLabel)
	b.WriteString("\n")

	textareaStyle := inputStyle
	if e.focused == int(FieldBody) {
		textareaStyle = focusedInputStyle
	}
	b.WriteString(textareaStyle.Render(e.body.View()))
	b.WriteString("\n\n")

	if e.draft != nil && len(e.draft.AttachmentPaths) > 0 {
		attachmentStyle := lipgloss.NewStyle().
			Foreground(Colors.Muted).
			Italic(true)
		attachments := strings.Join(e.draft.AttachmentPaths, ", ")
		b.WriteString(attachmentStyle.Render(fmt.Sprintf("Attachments: %s", attachments)))
		b.WriteString("\n")
	}

	if e.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(Colors.Error).
			Bold(true)
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", e.err)))
		b.WriteString("\n")
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(Colors.Muted).
		Padding(1, 0)

	help := helpStyle.Render(
		"Tab: next field • Shift+Tab: prev field • Ctrl+S: Send • Esc: Cancel",
	)
	b.WriteString(help)

	return b.String()
}

func (e *EmailComposer) Result() EmailComposerResult {
	return e.result
}

var EmailComposerKeys = emailComposerKeyMap{
	Send: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "send"),
	),
	SaveDraft: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "save draft"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
	Backtab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev field"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

type emailComposerKeyMap struct {
	Send      key.Binding
	SaveDraft key.Binding
	Cancel    key.Binding
	Tab       key.Binding
	Backtab   key.Binding
	Quit      key.Binding
}

func (k emailComposerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Send, k.SaveDraft, k.Cancel}
}

func (k emailComposerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.Backtab},
		{k.Send, k.SaveDraft, k.Cancel},
	}
}

type EmailComposerResultMsg struct {
	Result EmailComposerResult
}

type EmailComposerOpenMsg struct {
	Draft   *apply.EmailDraft
	Job     *job.Job
	Profile profile.Profile
}

func findPDF(texPath string) string {
	if texPath == "" {
		return ""
	}
	pdf := strings.TrimSuffix(texPath, ".tex") + ".pdf"
	if _, err := os.Stat(pdf); err == nil {
		return pdf
	}
	return ""
}

func RunEmailComposer(draft *apply.EmailDraft) (*apply.EmailDraft, EmailAction, error) {
	composer := NewEmailComposer(nil, profile.Profile{}, "")
	composer.SetDraft(draft)

	p := tea.NewProgram(composer)
	final, err := p.Run()
	if err != nil {
		return nil, ActionCancel, err
	}

	result := final.(*EmailComposer).Result()
	if result.Cancelled {
		return nil, ActionCancel, nil
	}

	return result.Draft, result.Action, nil
}
