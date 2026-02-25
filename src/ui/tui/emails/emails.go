package emails

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/theme"
)

// Email represents a single email draft or sent message.
type Email struct {
	Status  string // "draft" | "sent"
	Company string
	Subject string
	Date    string
	JobID   string
	Body    string
	From    string
	To      string
	Attach  string
}

// Model holds email list state.
type Model struct {
	Emails        []Email
	SelectedIndex int
	Width         int
	Height        int
}

// New creates a new emails model.
func New(emails []Email, w, h int) Model {
	return Model{
		Emails:        emails,
		SelectedIndex: 0,
		Width:         w,
		Height:        h,
	}
}

// OpenComposeMsg signals opening compose for a specific email.
type OpenComposeMsg struct{ Index int }

// GenerateMsg signals new email generation.
type GenerateMsg struct{}

// DeleteMsg signals email deletion.
type DeleteMsg struct{ Index int }

// BackMsg signals going back.
type BackMsg struct{}

// Update handles email list key events.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.SelectedIndex < len(m.Emails)-1 {
				m.SelectedIndex++
			}
		case "k", "up":
			if m.SelectedIndex > 0 {
				m.SelectedIndex--
			}
		case "enter":
			if len(m.Emails) > 0 {
				return m, func() tea.Msg { return OpenComposeMsg{Index: m.SelectedIndex} }
			}
		case "g":
			return m, func() tea.Msg { return GenerateMsg{} }
		case "d":
			if len(m.Emails) > 0 {
				return m, func() tea.Msg { return DeleteMsg{Index: m.SelectedIndex} }
			}
		case "esc":
			return m, func() tea.Msg { return BackMsg{} }
		}
	}
	return m, nil
}

// DraftCount returns the number of draft emails.
func (m Model) DraftCount() int {
	count := 0
	for _, e := range m.Emails {
		if e.Status == "draft" {
			count++
		}
	}
	return count
}

// SentCount returns the number of sent emails.
func (m Model) SentCount() int {
	count := 0
	for _, e := range m.Emails {
		if e.Status == "sent" {
			count++
		}
	}
	return count
}

// View renders the email list.
func (m Model) View() string {
	contentH := m.Height - 2
	if contentH < 1 {
		contentH = 1
	}

	bg := lipgloss.NewStyle().Background(theme.Background)

	// Column headers
	statusW := 8
	companyW := m.Width * 18 / 100
	if companyW < 10 {
		companyW = 10
	}
	dateW := 10
	subjectW := m.Width - statusW - companyW - dateW - 8

	header := theme.EmailHeaderStyle.Render(
		padRight("STATUS", statusW) +
			padRight("COMPANY", companyW) +
			padRight("SUBJECT", subjectW) +
			padRight("DATE", dateW))
	headerBorder := bg.Foreground(theme.BorderColor).Width(m.Width).Render(strings.Repeat("-", m.Width-4))

	var rows []string
	rows = append(rows, lipgloss.NewStyle().Background(theme.Background).Width(m.Width).PaddingLeft(2).Render(header))
	rows = append(rows, headerBorder)

	for i, e := range m.Emails {
		selected := i == m.SelectedIndex

		// Status badge
		var statusStr string
		if e.Status == "draft" {
			statusStr = theme.EmailStatusDraftStyle.Render("● draft")
		} else {
			statusStr = theme.EmailStatusSentStyle.Render("✓ sent ")
		}

		// Company
		company := e.Company
		if len(company) > companyW-2 {
			company = company[:companyW-5] + "..."
		}

		// Subject
		subject := e.Subject
		if len(subject) > subjectW-2 {
			subject = subject[:subjectW-5] + "..."
		}

		var companyStyle, subjectStyle lipgloss.Style
		var rowBg lipgloss.Color
		if selected {
			rowBg = lipgloss.Color("#091f28")
			companyStyle = theme.EmailCompanySelectedStyle.Background(rowBg)
			subjectStyle = theme.EmailSubjectStyle.Background(rowBg)
		} else {
			rowBg = theme.Background
			companyStyle = theme.EmailCompanyStyle.Background(rowBg)
			subjectStyle = theme.EmailSubjectStyle.Background(rowBg)
		}
		dateStyle := theme.EmailDateStyle.Background(rowBg)

		content := lipgloss.NewStyle().Background(rowBg).Render(statusStr) +
			lipgloss.NewStyle().Background(rowBg).Render(" ") +
			companyStyle.Width(companyW).Render(company) +
			subjectStyle.Width(subjectW).Render(subject) +
			dateStyle.Width(dateW).Render(e.Date)

		var row string
		if selected {
			row = lipgloss.NewStyle().
				Background(rowBg).
				Width(m.Width).
				BorderStyle(lipgloss.NormalBorder()).
				BorderLeft(true).
				BorderForeground(theme.Cyan).
				PaddingLeft(1).
				Render(content)
		} else {
			row = lipgloss.NewStyle().
				Background(rowBg).
				Width(m.Width).
				PaddingLeft(2).
				Render(content)
		}
		rows = append(rows, row)
	}

	body := lipgloss.JoinVertical(lipgloss.Left, rows...)
	bodyH := lipgloss.Height(body)
	if bodyH < contentH {
		pad := bg.Width(m.Width).Height(contentH - bodyH).Render("")
		body = lipgloss.JoinVertical(lipgloss.Left, body, pad)
	}

	return body
}

// TopBarRight returns "Drafts: N · Sent: N".
func (m Model) TopBarRight() string {
	return fmt.Sprintf("Drafts: %d  ·  Sent: %d", m.DraftCount(), m.SentCount())
}

func padRight(s string, w int) string {
	if len(s) >= w {
		return s[:w]
	}
	return s + strings.Repeat(" ", w-len(s))
}
