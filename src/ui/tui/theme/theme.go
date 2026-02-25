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

	// ── Filter fields ──────────────────────────────────────
	FieldRowStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Text)

	FieldRowFocusedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#13102a")).
				Foreground(Bright)

	FieldLabelStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Subtle)

	FieldLabelFocusedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#13102a")).
				Foreground(Purple)

	CursorStyle = lipgloss.NewStyle().
			Background(Purple).
			Foreground(Purple)

	// ── Profiles ───────────────────────────────────────────
	ProfileItemStyle = lipgloss.NewStyle().
				Background(Background).
				Foreground(Text).
				BorderStyle(lipgloss.NormalBorder()).
				BorderLeft(true).
				BorderForeground(Background)

	ProfileItemSelectedStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("#13102a")).
					Foreground(Purple).
					BorderStyle(lipgloss.NormalBorder()).
					BorderLeft(true).
					BorderForeground(Purple)

	ActiveBadgeStyle = lipgloss.NewStyle().
				Foreground(Green).
				Background(lipgloss.Color("#0d2b1e"))

	DetailSectionTitleStyle = lipgloss.NewStyle().
				Foreground(Muted).
				Bold(true)

	DetailKeyStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	DetailValStyle = lipgloss.NewStyle().
			Foreground(Text)

	DetailValYesStyle = lipgloss.NewStyle().
				Foreground(Green)

	DetailValNoStyle = lipgloss.NewStyle().
				Foreground(Muted)

	DetailValAccentStyle = lipgloss.NewStyle().
				Foreground(Cyan)

	BarLabelStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	BarTrackStyle = lipgloss.NewStyle().
			Foreground(Dim)

	BarFillStyle = lipgloss.NewStyle().
			Foreground(Purple)

	BarPctStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	// ── Help ───────────────────────────────────────────────
	HelpSectionTitleStyle = lipgloss.NewStyle().
				Foreground(Cyan).
				Bold(true)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Background(lipgloss.Color("#0d2b33"))

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(Text)

	// ── Scraping ───────────────────────────────────────────
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(Purple)

	ProgressFillStyle = lipgloss.NewStyle().
				Foreground(Purple)

	ProgressTrackStyle = lipgloss.NewStyle().
				Foreground(Dim)

	ProgressPctStyle = lipgloss.NewStyle().
				Foreground(Purple).
				Bold(true)

	// ── Emails ─────────────────────────────────────────────
	EmailRowStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Text)

	EmailRowSelectedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#091f28")).
				Foreground(Bright).
				BorderStyle(lipgloss.NormalBorder()).
				BorderLeft(true).
				BorderForeground(Cyan)

	EmailStatusDraftStyle = lipgloss.NewStyle().
				Foreground(Yellow)

	EmailStatusSentStyle = lipgloss.NewStyle().
				Foreground(Green)

	EmailCompanyStyle = lipgloss.NewStyle().
				Foreground(Bright)

	EmailCompanySelectedStyle = lipgloss.NewStyle().
					Foreground(Cyan)

	EmailSubjectStyle = lipgloss.NewStyle().
				Foreground(Subtle)

	EmailDateStyle = lipgloss.NewStyle().
			Foreground(Muted)

	EmailHeaderStyle = lipgloss.NewStyle().
				Foreground(Muted)

	// ── Compose ────────────────────────────────────────────
	PopFieldStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Text)

	PopFieldActiveStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#13102a")).
				Foreground(Bright)

	PopLabelStyle = lipgloss.NewStyle().
			Foreground(Muted)

	PopLabelActiveStyle = lipgloss.NewStyle().
				Foreground(Purple)

	PopGenHintStyle = lipgloss.NewStyle().
			Foreground(Purple)

	PopBodyStyle = lipgloss.NewStyle().
			Foreground(Text)

	// ── CV Wizard ──────────────────────────────────────────
	CVBreadcrumbStyle = lipgloss.NewStyle().
				Foreground(Muted)

	CVBreadcrumbCurrentStyle = lipgloss.NewStyle().
					Foreground(Bright).
					Bold(true).
					Background(Surface4)

	CVBreadcrumbDoneStyle = lipgloss.NewStyle().
				Foreground(Purple)

	CVBreadcrumbSepStyle = lipgloss.NewStyle().
				Foreground(Dim)

	CVChoiceStyle = lipgloss.NewStyle().
			Foreground(Text).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(BorderColor)

	CVChoiceSelectedStyle = lipgloss.NewStyle().
				Foreground(Bright).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(Purple).
				Background(lipgloss.Color("#13102a"))

	CVChoiceIconStyle = lipgloss.NewStyle().
				Foreground(Muted)

	CVChoiceIconSelectedStyle = lipgloss.NewStyle().
					Foreground(Purple)

	CVFormFieldStyle = lipgloss.NewStyle().
				Background(Background).
				Foreground(Text).
				BorderStyle(lipgloss.NormalBorder()).
				BorderLeft(true).
				BorderForeground(Background)

	CVFormFieldFocusedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#13102a")).
				Foreground(Bright).
				BorderStyle(lipgloss.NormalBorder()).
				BorderLeft(true).
				BorderForeground(Purple)

	CVFormLabelStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	CVFormLabelFocusedStyle = lipgloss.NewStyle().
				Foreground(Purple)

	CVTagStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Background(lipgloss.Color("#0d2b33"))

	CVTagAddStyle = lipgloss.NewStyle().
			Foreground(Muted)

	CVTagAddFocusedStyle = lipgloss.NewStyle().
				Foreground(Purple)

	CVTextareaStyle = lipgloss.NewStyle().
			Foreground(Text).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#2d2157"))

	CVReviewKeyStyle = lipgloss.NewStyle().
				Foreground(Muted)

	CVReviewValStyle = lipgloss.NewStyle().
				Foreground(Text)

	CVStepDotActiveStyle = lipgloss.NewStyle().
				Foreground(Purple)

	CVStepDotInactiveStyle = lipgloss.NewStyle().
				Foreground(Dim)

	// ── Placeholder text ───────────────────────────────────
	PlaceholderStyle = lipgloss.NewStyle().
			Foreground(Muted)

	// ── Section label ──────────────────────────────────────
	SectionLabelStyle = lipgloss.NewStyle().
				Foreground(Muted).
				Bold(true)
)
