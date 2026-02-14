package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Catppuccin-inspired palette
	colorRosewater = lipgloss.Color("#f5e0dc")
	colorFlamingo  = lipgloss.Color("#f2cdcd")
	colorPink      = lipgloss.Color("#f5c2e7")
	colorMauve     = lipgloss.Color("#cba6f7")
	colorRed       = lipgloss.Color("#f38ba8")
	colorMaroon    = lipgloss.Color("#eba0ac")
	colorPeach     = lipgloss.Color("#fab387")
	colorYellow    = lipgloss.Color("#f9e2af")
	colorGreen     = lipgloss.Color("#a6e3a1")
	colorTeal      = lipgloss.Color("#94e2d5")
	colorSky       = lipgloss.Color("#89dceb")
	colorSapphire  = lipgloss.Color("#74c7ec")
	colorBlue      = lipgloss.Color("#89b4fa")
	colorLavender  = lipgloss.Color("#b4befe")
	colorText      = lipgloss.Color("#cdd6f4")
	colorSubtext1  = lipgloss.Color("#bac2de")
	colorSubtext0  = lipgloss.Color("#a6adc8")
	colorOverlay2  = lipgloss.Color("#9399b2")
	colorOverlay1  = lipgloss.Color("#7f849c")
	colorOverlay0  = lipgloss.Color("#6c7086")
	colorSurface2  = lipgloss.Color("#585b70")
	colorSurface1  = lipgloss.Color("#45475a")
	colorSurface0  = lipgloss.Color("#313244")
	colorBase      = lipgloss.Color("#1e1e2e")
	colorMantle    = lipgloss.Color("#181825")
	colorCrust     = lipgloss.Color("#11111b")

	// Styles
	styleTitle = lipgloss.NewStyle().
			Foreground(colorMauve).
			Bold(true).
			Padding(0, 1)

	styleSelected = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	styleNormal = lipgloss.NewStyle().
			Foreground(colorText)

	styleDim = lipgloss.NewStyle().
			Foreground(colorOverlay1)

	styleActiveTab = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMauve).
			Padding(0, 1).
			Foreground(colorMauve)

	styleInactiveTab = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSurface1).
			Padding(0, 1).
			Foreground(colorOverlay1)

	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSurface1).
			Padding(1).
			Margin(0, 1)

	queryBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMauve).
			Padding(1).
			Margin(0, 1)

	styleFooter = lipgloss.NewStyle().
			Foreground(colorSubtext1).
			Background(colorSurface0).
			Padding(0, 1)

	styleError = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	styleSuccess = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)
)
