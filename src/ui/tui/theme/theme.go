package theme

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	Background  = lipgloss.Color("#0e0e0e")
	Surface     = lipgloss.Color("#161616")
	Surface2    = lipgloss.Color("#1c1c1c")
	Surface3    = lipgloss.Color("#262626")
	Surface4    = lipgloss.Color("#333333")
	BorderColor = lipgloss.Color("#2a2a2a")
	Muted       = lipgloss.Color("#444444")
	Subtle      = lipgloss.Color("#686868")
	Text        = lipgloss.Color("#c8c8c8")
	Bright      = lipgloss.Color("#f0f0f0")
	Dim         = lipgloss.Color("#3a3a3a")
	Yellow      = lipgloss.Color("#f0c060")
	Green       = lipgloss.Color("#50e3a4")
	Cyan        = lipgloss.Color("#4cc9f0")
	Purple      = lipgloss.Color("#a78bfa")
	Accent      = lipgloss.Color("#7b61ff")
)

var (
	BaseStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Text)

	// TopBarStyle is applied only via renderTopBar() which builds its own
	// row — this style is kept for any single-shot callers.
	TopBarStyle = lipgloss.NewStyle().
			Background(Surface).
			Foreground(Text).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(BorderColor).
			BorderBackground(Surface)

	StatusBarStyle = lipgloss.NewStyle().
			Background(Surface).
			Foreground(Text)

	StatusLabelStyle = lipgloss.NewStyle().
				Background(Surface).
				Foreground(Subtle)

	ContentStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Text)

	EmptyStateStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Muted)

	EmptyASCIIStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Dim)

	EmptyHeadlineStyle = lipgloss.NewStyle().
				Background(Background).
				Foreground(Bright).
				Bold(true)

	EmptySubStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Subtle)

	JobItemStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Text)

	JobItemSelectedStyle = lipgloss.NewStyle().
				Background(Surface3).
				Foreground(Bright)

	JobScoreStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Cyan).
			Bold(true)

	JobTrapsStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Yellow).
			Bold(true)

	JobCompanyStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Subtle)

	JobSourceStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Muted)

	// KbdStyle — keycap badge: cyan text on a subtle cyan-tinted background.
	KbdStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0d2b33")).
			Foreground(Cyan).
			PaddingLeft(1).
			PaddingRight(1)

	SepStyle = lipgloss.NewStyle().
			Background(Surface).
			Foreground(Dim)

	ProgressStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Purple)

	SuccessStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Green)

	WarningStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Yellow)

	ErrorStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(lipgloss.Color("#ff5555"))

	ModalTopBarStyle = lipgloss.NewStyle().
				Background(Surface2).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(BorderColor).
				BorderBackground(Surface2).
				Height(1)

	ModalTitleStyle = lipgloss.NewStyle().
			Background(Surface2).
			Foreground(Bright).
			Bold(true)

	ModalHintStyle = lipgloss.NewStyle().
			Background(Surface2).
			Foreground(Muted)
)