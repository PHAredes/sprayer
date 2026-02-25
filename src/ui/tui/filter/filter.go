package filter

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sprayer/src/ui/tui/theme"
)

// Field holds a single editable filter field.
type Field struct {
	Label string
	Value string
}

// Model holds filter screen state.
type Model struct {
	Fields     []Field
	FocusIndex int
	Width      int
	Height     int
}

// New creates a filter model pre-populated with the given profile values.
func New(keywords, exclude, locations, companies, minScore string, w, h int) Model {
	return Model{
		Fields: []Field{
			{Label: "Keywords", Value: keywords},
			{Label: "Exclude", Value: exclude},
			{Label: "Locations", Value: locations},
			{Label: "Companies", Value: companies},
			{Label: "Min Score", Value: minScore},
		},
		FocusIndex: 0,
		Width:      w,
		Height:     h,
	}
}

// FilterApplyMsg carries updated filter values.
type FilterApplyMsg struct {
	Keywords  string
	Exclude   string
	Locations string
	Companies string
	MinScore  string
}

// FilterCancelMsg signals cancellation.
type FilterCancelMsg struct{}

// Update handles key events inside the filter screen.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.FocusIndex = (m.FocusIndex + 1) % len(m.Fields)
		case "shift+tab", "up":
			m.FocusIndex = (m.FocusIndex - 1 + len(m.Fields)) % len(m.Fields)
		case "enter":
			return m, func() tea.Msg {
				return FilterApplyMsg{
					Keywords:  m.Fields[0].Value,
					Exclude:   m.Fields[1].Value,
					Locations: m.Fields[2].Value,
					Companies: m.Fields[3].Value,
					MinScore:  m.Fields[4].Value,
				}
			}
		case "esc":
			return m, func() tea.Msg { return FilterCancelMsg{} }
		case "backspace":
			f := &m.Fields[m.FocusIndex]
			if len(f.Value) > 0 {
				f.Value = f.Value[:len(f.Value)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.Fields[m.FocusIndex].Value += msg.String()
			} else if msg.String() == " " {
				m.Fields[m.FocusIndex].Value += " "
			}
		}
	}
	return m, nil
}

// View renders the filter form.
func (m Model) View() string {
	contentH := m.Height - 4 // topbar + modal-topbar + statusbar + border-ish
	if contentH < 1 {
		contentH = 1
	}

	labelW := 12

	var rows []string
	for i, f := range m.Fields {
		focused := i == m.FocusIndex

		var labelStyle, valStyle lipgloss.Style
		var rowBg lipgloss.Color
		if focused {
			labelStyle = theme.FieldLabelFocusedStyle
			valStyle = lipgloss.NewStyle().Foreground(theme.Bright)
			rowBg = lipgloss.Color("#13102a")
		} else {
			labelStyle = theme.FieldLabelStyle
			valStyle = lipgloss.NewStyle().Foreground(theme.Text)
			rowBg = theme.Background
		}

		label := labelStyle.Width(labelW).Align(lipgloss.Right).Render(f.Label)
		sep := lipgloss.NewStyle().Background(rowBg).Render("  ")
		val := valStyle.Background(rowBg).Render(f.Value)
		cursor := ""
		if focused {
			cursor = theme.CursorStyle.Render("|")
		}

		rowContent := label + sep + val + cursor
		row := lipgloss.NewStyle().
			Background(rowBg).
			Width(m.Width - 4).
			PaddingLeft(1).
			PaddingRight(1).
			Render(rowContent)
		rows = append(rows, row)
	}

	// Pad remaining space
	body := lipgloss.JoinVertical(lipgloss.Left, rows...)
	bodyH := lipgloss.Height(body)
	if bodyH < contentH {
		pad := lipgloss.NewStyle().
			Background(theme.Background).
			Width(m.Width).
			Height(contentH - bodyH).
			Render("")
		body = lipgloss.JoinVertical(lipgloss.Left, body, pad)
	}

	return lipgloss.NewStyle().
		Background(theme.Background).
		Width(m.Width).
		Padding(1, 2).
		Render(body)
}

// ModalTopBar renders the filter modal top bar.
func (m Model) ModalTopBar() string {
	title := theme.ModalTitleStyle.Render("Filter Configuration")
	hint := theme.ModalHintStyle.Render("tab next  |  enter apply  |  esc cancel")

	titleW := lipgloss.Width(title)
	hintW := lipgloss.Width(hint)
	gap := m.Width - titleW - hintW - 4
	if gap < 1 {
		gap = 1
	}
	pad := lipgloss.NewStyle().Background(theme.Surface2).Width(gap).Render("")

	row := title + pad + hint
	return theme.ModalTopBarStyle.Width(m.Width).PaddingLeft(2).PaddingRight(2).Render(row)
}

// StatusBarRight returns the right side text for the filter status bar.
func StatusBarRight(profileName string) string {
	return "Profile: " + profileName + "  |  " + strings.Repeat("", 0) + "No traps"
}
