package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette inspired by opencode and terminal.shop
// Clean, minimal, high-contrast dark theme
var Colors = struct {
	// Core palette
	Base    lipgloss.Color
	Surface lipgloss.Color
	Overlay lipgloss.Color
	Muted   lipgloss.Color
	Subtle  lipgloss.Color
	Text    lipgloss.Color

	// Accent colors (used sparingly)
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color

	// Semantic colors
	Success lipgloss.Color
	Warning lipgloss.Color
	Error   lipgloss.Color
	Info    lipgloss.Color

	// Interactive
	Selection lipgloss.Color
	Highlight lipgloss.Color
	Border    lipgloss.Color
	BorderDim lipgloss.Color
}{
	// Core - minimal, almost monochrome
	Base:    lipgloss.Color("#0a0a0a"), // Near black
	Surface: lipgloss.Color("#111111"), // Slight lift
	Overlay: lipgloss.Color("#1a1a1a"), // Cards, elevated elements
	Muted:   lipgloss.Color("#404040"), // Dimmed text
	Subtle:  lipgloss.Color("#666666"), // Secondary text
	Text:    lipgloss.Color("#e4e4e4"), // Primary text - high contrast

	// Accent - subtle, professional
	Primary:   lipgloss.Color("#888888"), // Muted gray for primary actions
	Secondary: lipgloss.Color("#555555"), // Darker gray
	Accent:    lipgloss.Color("#c0c0c0"), // Bright gray for highlights

	// Semantic - desaturated
	Success: lipgloss.Color("#4a9a6a"), // Muted green
	Warning: lipgloss.Color("#a08040"), // Muted amber
	Error:   lipgloss.Color("#a05050"), // Muted red
	Info:    lipgloss.Color("#5080a0"), // Muted blue

	// Interactive
	Selection: lipgloss.Color("#1e1e1e"), // Selection background
	Highlight: lipgloss.Color("#ffffff"), // Bright highlights
	Border:    lipgloss.Color("#222222"), // Visible borders
	BorderDim: lipgloss.Color("#1a1a1a"), // Subtle borders
}

// Styles - minimal, clean, functional
var Styles = struct {
	// Layout
	App       lipgloss.Style
	Container lipgloss.Style
	Card      lipgloss.Style
	Content   lipgloss.Style
	Detail    lipgloss.Style

	// Typography
	Title      lipgloss.Style
	Header     lipgloss.Style
	HeaderText lipgloss.Style
	Subtitle   lipgloss.Style
	Text       lipgloss.Style
	Muted      lipgloss.Style
	MutedText  lipgloss.Style
	Bold       lipgloss.Style

	// Semantic text
	Success     lipgloss.Style
	SuccessText lipgloss.Style
	Warning     lipgloss.Style
	WarningText lipgloss.Style
	Error       lipgloss.Style
	ErrorText   lipgloss.Style
	Info        lipgloss.Style

	// List styles
	List         lipgloss.Style
	ListItem     lipgloss.Style
	ListTitle    lipgloss.Style
	SelectedItem lipgloss.Style
	ActiveItem   lipgloss.Style

	// Detail styles
	DetailTitle   lipgloss.Style
	DetailSection lipgloss.Style

	// Inputs
	Input        lipgloss.Style
	InputFocused lipgloss.Style
	InputPrompt  lipgloss.Style

	// Buttons
	Button       lipgloss.Style
	ButtonActive lipgloss.Style

	// Status
	StatusBar  lipgloss.Style
	StatusText lipgloss.Style

	// Misc
	Border      lipgloss.Style
	BorderLight lipgloss.Style
	Help        lipgloss.Style
	HelpBox     lipgloss.Style
	Loading     lipgloss.Style
	Scraping    lipgloss.Style
	FilterBox   lipgloss.Style
}{
	// Layout
	App: lipgloss.NewStyle().
		Background(Colors.Base).
		Foreground(Colors.Text),

	Container: lipgloss.NewStyle().
		Padding(0, 1),

	Card: lipgloss.NewStyle().
		Background(Colors.Overlay).
		BorderLeft(true).
		BorderForeground(Colors.Border).
		BorderLeftForeground(Colors.Primary).
		Padding(0, 1).
		MarginLeft(1),

	Content: lipgloss.NewStyle().
		Foreground(Colors.Text).
		Padding(0, 1),

	Detail: lipgloss.NewStyle().
		Background(Colors.Base).
		Padding(1, 2),

	// Typography
	Title: lipgloss.NewStyle().
		Foreground(Colors.Text).
		Bold(true).
		MarginBottom(1),

	Header: lipgloss.NewStyle().
		Foreground(Colors.Text).
		Bold(true).
		Padding(0, 1),

	HeaderText: lipgloss.NewStyle().
		Foreground(Colors.Subtle).
		Bold(false),

	Subtitle: lipgloss.NewStyle().
		Foreground(Colors.Subtle).
		Faint(true),

	Text: lipgloss.NewStyle().
		Foreground(Colors.Text),

	Muted: lipgloss.NewStyle().
		Foreground(Colors.Muted),

	MutedText: lipgloss.NewStyle().
		Foreground(Colors.Muted).
		Faint(true),

	Bold: lipgloss.NewStyle().
		Foreground(Colors.Text).
		Bold(true),

	// Semantic text
	Success: lipgloss.NewStyle().
		Foreground(Colors.Success),

	SuccessText: lipgloss.NewStyle().
		Foreground(Colors.Success).
		Bold(true),

	Warning: lipgloss.NewStyle().
		Foreground(Colors.Warning),

	WarningText: lipgloss.NewStyle().
		Foreground(Colors.Warning).
		Bold(true),

	Error: lipgloss.NewStyle().
		Foreground(Colors.Error),

	ErrorText: lipgloss.NewStyle().
		Foreground(Colors.Error).
		Bold(true),

	Info: lipgloss.NewStyle().
		Foreground(Colors.Info),

	// List styles
	List: lipgloss.NewStyle().
		Background(Colors.Base),

	ListItem: lipgloss.NewStyle().
		Foreground(Colors.Text).
		Padding(0, 1),

	ListTitle: lipgloss.NewStyle().
		Foreground(Colors.Text).
		Bold(true).
		MarginBottom(1),

	SelectedItem: lipgloss.NewStyle().
		Background(Colors.Selection).
		Foreground(Colors.Text).
		Bold(true).
		Padding(0, 1),

	ActiveItem: lipgloss.NewStyle().
		Foreground(Colors.Highlight).
		Bold(true).
		Padding(0, 1),

	// Detail styles
	DetailTitle: lipgloss.NewStyle().
		Foreground(Colors.Text).
		Bold(true).
		MarginBottom(1),

	DetailSection: lipgloss.NewStyle().
		MarginBottom(1),

	// Inputs
	Input: lipgloss.NewStyle().
		Background(Colors.Surface).
		Foreground(Colors.Text).
		Padding(0, 1),

	InputFocused: lipgloss.NewStyle().
		Background(Colors.Surface).
		Foreground(Colors.Highlight).
		Padding(0, 1),

	InputPrompt: lipgloss.NewStyle().
		Foreground(Colors.Subtle),

	// Buttons
	Button: lipgloss.NewStyle().
		Background(Colors.Overlay).
		Foreground(Colors.Text).
		Padding(0, 2),

	ButtonActive: lipgloss.NewStyle().
		Background(Colors.Selection).
		Foreground(Colors.Highlight).
		Padding(0, 2).
		Bold(true),

	// Status
	StatusBar: lipgloss.NewStyle().
		Background(Colors.Surface).
		Foreground(Colors.Subtle).
		Height(1).
		Padding(0, 1),

	StatusText: lipgloss.NewStyle().
		Foreground(Colors.Subtle),

	// Misc
	Border: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.BorderDim),

	BorderLight: lipgloss.NewStyle().
		BorderLeft(true).
		BorderForeground(Colors.Border),

	Help: lipgloss.NewStyle().
		Background(Colors.Overlay).
		Foreground(Colors.Text).
		Padding(1, 2),

	HelpBox: lipgloss.NewStyle().
		Background(Colors.Overlay).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Border).
		Padding(1, 2),

	Loading: lipgloss.NewStyle().
		Foreground(Colors.Subtle).
		Padding(2),

	Scraping: lipgloss.NewStyle().
		Foreground(Colors.Subtle).
		Padding(2).
		Align(lipgloss.Center),

	FilterBox: lipgloss.NewStyle().
		Background(Colors.Surface).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Border).
		Padding(1, 2),
}

// Layout utilities
var Layout = struct {
	HorizontalJoin func(...string) string
	VerticalJoin   func(...string) string
}{
	HorizontalJoin: func(parts ...string) string {
		return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	},
	VerticalJoin: func(parts ...string) string {
		return lipgloss.JoinVertical(lipgloss.Left, parts...)
	},
}

// Text rendering helpers
func MutedText(s string) string {
	return Styles.MutedText.Render(s)
}

func HighlightText(s string) string {
	return Styles.ActiveItem.Render(s)
}

func ErrorText(s string) string {
	return Styles.ErrorText.Render(s)
}

func SuccessText(s string) string {
	return Styles.SuccessText.Render(s)
}

func WarningText(s string) string {
	return Styles.WarningText.Render(s)
}

func InfoText(s string) string {
	return Styles.Info.Render(s)
}

// Indent helper
func Indent(s string, n int) string {
	return lipgloss.NewStyle().PaddingLeft(n).Render(s)
}

// Truncate with ellipsis
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
