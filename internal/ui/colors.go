package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// CHARM color palette inspired by terminal.shop and CHARM ecosystem
var Colors = struct {
	// Primary colors
	Primary lipgloss.Color
	Accent  lipgloss.Color
	Success lipgloss.Color
	Warning lipgloss.Color
	Error   lipgloss.Color

	// Neutral colors
	Text   lipgloss.Color
	Subtle lipgloss.Color
	Muted  lipgloss.Color

	// Surface colors
	Background lipgloss.Color
	Surface    lipgloss.Color
	Elevated   lipgloss.Color
	Border     lipgloss.Color

	// Interactive colors
	Selected lipgloss.Color
	Hover    lipgloss.Color
	Active   lipgloss.Color

	// Special colors
	Highlight lipgloss.Color
	Dim       lipgloss.Color
}{
	// Primary colors - following CHARM's signature purple/blue theme
	Primary: lipgloss.Color("#7D56F4"), // CHARM purple
	Accent:  lipgloss.Color("#00A3FF"), // Bright blue
	Success: lipgloss.Color("#04B575"), // Green
	Warning: lipgloss.Color("#F5A524"), // Orange
	Error:   lipgloss.Color("#E54D2E"), // Red

	// Neutral colors
	Text:   lipgloss.Color("#E6E1E5"), // Light gray
	Subtle: lipgloss.Color("#C9C5CA"), // Medium gray
	Muted:  lipgloss.Color("#9F9BA1"), // Dark gray

	// Surface colors
	Background: lipgloss.Color("#161618"), // Dark background
	Surface:    lipgloss.Color("#1C1C1F"), // Slightly lighter
	Elevated:   lipgloss.Color("#232326"), // Card background
	Border:     lipgloss.Color("#2F2F32"), // Border color

	// Interactive colors
	Selected: lipgloss.Color("#3A3A3D"), // Selection background
	Hover:    lipgloss.Color("#353538"), // Hover background
	Active:   lipgloss.Color("#7D56F4"), // Active accent

	// Special colors
	Highlight: lipgloss.Color("#FFD60A"), // Yellow highlight
	Dim:       lipgloss.Color("#706F74"), // Dimmed text
}

// Styles contains all the styled components following CHARM design patterns
var Styles = struct {
	// Layout styles
	App       lipgloss.Style
	Container lipgloss.Style
	Card      lipgloss.Style

	// Header styles
	Header     lipgloss.Style
	HeaderText lipgloss.Style
	Title      lipgloss.Style
	Subtitle   lipgloss.Style

	// Content styles
	Content     lipgloss.Style
	Text        lipgloss.Style
	MutedText   lipgloss.Style
	ErrorText   lipgloss.Style
	SuccessText lipgloss.Style

	// List styles
	List         lipgloss.Style
	ListItem     lipgloss.Style
	SelectedItem lipgloss.Style
	ListTitle    lipgloss.Style

	// Detail styles
	Detail        lipgloss.Style
	DetailTitle   lipgloss.Style
	DetailSection lipgloss.Style

	// Input styles
	Input        lipgloss.Style
	InputFocused lipgloss.Style
	InputPrompt  lipgloss.Style

	// Button styles
	Button        lipgloss.Style
	ButtonFocused lipgloss.Style
	ButtonActive  lipgloss.Style

	// Status styles
	StatusBar     lipgloss.Style
	StatusText    lipgloss.Style
	StatusSuccess lipgloss.Style
	StatusError   lipgloss.Style

	// Special styles
	Loading   lipgloss.Style
	Scraping  lipgloss.Style
	HelpBox   lipgloss.Style
	FilterBox lipgloss.Style

	// Border styles
	Border       lipgloss.Style
	BorderFocus  lipgloss.Style
	BorderActive lipgloss.Style
}{
	// Layout styles
	App: lipgloss.NewStyle().
		Background(Colors.Background).
		Foreground(Colors.Text),

	Container: lipgloss.NewStyle().
		Padding(1, 2),

	Card: lipgloss.NewStyle().
		Background(Colors.Elevated).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Border).
		Padding(1, 2).
		Margin(0, 1),

	// Header styles
	Header: lipgloss.NewStyle().
		Background(Colors.Surface).
		Foreground(Colors.Text).
		Padding(0, 1).
		Height(1),

	HeaderText: lipgloss.NewStyle().
		Foreground(Colors.Subtle).
		Bold(false),

	Title: lipgloss.NewStyle().
		Foreground(Colors.Primary).
		Bold(true).
		Padding(0, 1),

	Subtitle: lipgloss.NewStyle().
		Foreground(Colors.Muted).
		Faint(true),

	// Content styles
	Content: lipgloss.NewStyle().
		Foreground(Colors.Text).
		Padding(0, 1),

	Text: lipgloss.NewStyle().
		Foreground(Colors.Text),

	MutedText: lipgloss.NewStyle().
		Foreground(Colors.Muted).
		Faint(true),

	ErrorText: lipgloss.NewStyle().
		Foreground(Colors.Error).
		Bold(true),

	SuccessText: lipgloss.NewStyle().
		Foreground(Colors.Success).
		Bold(true),

	// List styles
	List: lipgloss.NewStyle().
		Background(Colors.Background),

	ListItem: lipgloss.NewStyle().
		Foreground(Colors.Text).
		Padding(0, 1),

	SelectedItem: lipgloss.NewStyle().
		Background(Colors.Selected).
		Foreground(Colors.Accent).
		Bold(true).
		Padding(0, 1),

	ListTitle: lipgloss.NewStyle().
		Foreground(Colors.Primary).
		Bold(true).
		MarginBottom(1),

	// Detail styles
	Detail: lipgloss.NewStyle().
		Background(Colors.Background).
		Padding(1, 2),

	DetailTitle: lipgloss.NewStyle().
		Foreground(Colors.Primary).
		Bold(true).
		MarginBottom(1),

	DetailSection: lipgloss.NewStyle().
		MarginBottom(1),

	// Input styles
	Input: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Border).
		Padding(0, 1),

	InputFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Accent).
		Padding(0, 1),

	InputPrompt: lipgloss.NewStyle().
		Foreground(Colors.Subtle).
		MarginRight(1),

	// Button styles
	Button: lipgloss.NewStyle().
		Background(Colors.Surface).
		Foreground(Colors.Text).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Border).
		Padding(0, 2),

	ButtonFocused: lipgloss.NewStyle().
		Background(Colors.Hover).
		Foreground(Colors.Accent).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Accent).
		Padding(0, 2),

	ButtonActive: lipgloss.NewStyle().
		Background(Colors.Primary).
		Foreground(Colors.Text).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Primary).
		Padding(0, 2).
		Bold(true),

	// Status styles
	StatusBar: lipgloss.NewStyle().
		Background(Colors.Surface).
		Foreground(Colors.Text).
		Height(1).
		Padding(0, 1),

	StatusText: lipgloss.NewStyle().
		Foreground(Colors.Subtle),

	StatusSuccess: lipgloss.NewStyle().
		Foreground(Colors.Success),

	StatusError: lipgloss.NewStyle().
		Foreground(Colors.Error).
		Bold(true),

	// Special styles
	Loading: lipgloss.NewStyle().
		Foreground(Colors.Accent).
		Bold(true).
		Padding(2),

	Scraping: lipgloss.NewStyle().
		Foreground(Colors.Accent).
		Padding(2).
		Align(lipgloss.Center),

	HelpBox: lipgloss.NewStyle().
		Background(Colors.Elevated).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Border).
		Padding(1, 2),

	FilterBox: lipgloss.NewStyle().
		Background(Colors.Surface).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Border).
		Padding(1, 2),

	// Border styles
	Border: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Border),

	BorderFocus: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Accent),

	BorderActive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Primary),
}

// Layout provides layout utilities following CHARM patterns
var Layout = struct {
	// Common layouts
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
