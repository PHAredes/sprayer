package tui

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"sprayer/internal/job"
)

type ViewState int

const (
	EmptyState ViewState = iota
	JobList
	Filter
	Profiles
	Help
	Scraping
	Emails
	Compose
	CVEntry
	CVInfo
	CVSummary
	CVExperience
	CVSkills
	CVReview
)

type Model struct {
	jobs          []job.Job
	selectedIndex int
	profileName   string
	viewState     ViewState
	width         int
	height        int
}

func NewModel() Model {
	return Model{
		jobs:          []job.Job{},
		selectedIndex: 0,
		profileName:   "Default",
		viewState:     EmptyState,
		width:         80,
		height:        24,
	}
}

func (m *Model) SelectedIndex() int     { return m.selectedIndex }
func (m *Model) ViewState() ViewState   { return m.viewState }
func (m *Model) Jobs() []job.Job        { return m.jobs }
func (m *Model) SetJobs(jobs []job.Job) { m.jobs = jobs }

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "↓":
			if len(m.jobs) > 0 {
				m.selectedIndex = min(m.selectedIndex+1, len(m.jobs)-1)
				m.viewState = JobList
			}
		case "k", "↑":
			if len(m.jobs) > 0 {
				m.selectedIndex = max(m.selectedIndex-1, 0)
				m.viewState = JobList
			}
		case "s":
			m.viewState = Scraping
		case "f":
			m.viewState = Filter
		case "p":
			m.viewState = Profiles
		case "m":
			m.viewState = Emails
		case "a":
		case "?":
			m.viewState = Help
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderTopBar(),
		m.renderContent(),
		m.renderStatusBar(),
	)
}

// ── Top bar — single row ──────────────────────────────────────────────────────

func (m Model) renderTopBar() string {
	on := func(fg lipgloss.Color) lipgloss.Style {
		return lipgloss.NewStyle().Background(Surface).Foreground(fg)
	}

	left  := on(Subtle).Render("Profile: ") + on(Cyan).Render(m.profileName)
	title := on(Bright).Bold(true).Render("Sprayer")
	right := on(Subtle).Render("Jobs: ") + on(Yellow).Render(strconv.Itoa(len(m.jobs)))

	titleW := lipgloss.Width(title)
	sideW  := (m.width - titleW) / 2
	rightW := m.width - sideW - titleW

	leftBlock  := lipgloss.NewStyle().Background(Surface).Width(sideW).PaddingLeft(2).Render(left)
	rightBlock := lipgloss.NewStyle().Background(Surface).Width(rightW).PaddingRight(2).Align(lipgloss.Right).Render(right)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBlock, title, rightBlock)
}

// ── Content ───────────────────────────────────────────────────────────────────

func (m Model) renderContent() string {
	if len(m.jobs) == 0 {
		return m.renderEmptyState()
	}
	return m.renderJobList()
}

// contentHeight is height minus the two single-row bars.
func (m Model) contentHeight() int { return m.height - 2 }

func (m Model) renderEmptyState() string {
	availH := m.contentHeight()
	bg     := lipgloss.NewStyle().Background(Background)

	ascii := EmptyASCIIStyle.Render(
		"┌─────────────────────────┐\n" +
			"│  ░░░░░░░░░░░░░░░░░░░░░  │\n" +
			"│  ░  No jobs cached.  ░  │\n" +
			"│  ░░░░░░░░░░░░░░░░░░░░░  │\n" +
			"└─────────────────────────┘")

	headline := EmptyHeadlineStyle.Render("No jobs found")

	// Build each subtext line separately so we control their exact width.
	// Without this, JoinVertical pads shorter lines with spaces that inherit
	// the wrong (terminal default) background colour.
	subLines := m.emptySubLines()
	subW := 0
	for _, l := range subLines {
		if w := lipgloss.Width(l); w > subW {
			subW = w
		}
	}
	padded := make([]string, len(subLines))
	for i, l := range subLines {
		padded[i] = bg.Width(subW).Render(l)
	}
	subtext := lipgloss.JoinVertical(lipgloss.Left, padded...)

	sp     := bg.Render("")
	blockW := max(lipgloss.Width(ascii), lipgloss.Width(headline), subW)

	wrap := func(s string) string {
		return bg.Width(blockW).Align(lipgloss.Center).Render(s)
	}

	block := lipgloss.JoinVertical(lipgloss.Left,
		wrap(ascii),
		wrap(sp),
		wrap(headline),
		wrap(sp),
		wrap(subtext),
	)

	return lipgloss.Place(
		m.width, availH,
		lipgloss.Center, lipgloss.Center,
		block,
		lipgloss.WithWhitespaceBackground(Background),
	)
}

// emptySubLines returns the three hint lines as raw (unsized) strings.
func (m Model) emptySubLines() []string {
	prose := lipgloss.NewStyle().Background(Background).Foreground(Subtle)
	kbd   := lipgloss.NewStyle().Background(lipgloss.Color("#0d2b33")).Foreground(Cyan).
		PaddingLeft(1).PaddingRight(1)
	dot := prose.Render(" · ")

	return []string{
		prose.Render("Press ") + kbd.Render("s") + prose.Render(" to scrape from all configured sources."),
		prose.Render("Press ") + kbd.Render("f") + prose.Render(" to set filters") + dot + kbd.Render("p") + prose.Render(" to manage profiles."),
		prose.Render("Jobs appear here in real time as they're discovered."),
	}
}

func (m Model) renderJobList() string {
	availH := m.contentHeight()

	var lines []string
	for i, j := range m.jobs {
		var line string
		if i == m.selectedIndex {
			line = JobItemSelectedStyle.Width(m.width).Render(m.formatJobItem(j))
		} else {
			line = JobItemStyle.Width(m.width).Render(m.formatJobItem(j))
		}
		lines = append(lines, line)
	}
	for len(lines) < availH {
		lines = append(lines, ContentStyle.Width(m.width).Render(""))
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m Model) formatJobItem(j job.Job) string {
	scoreStr   := JobScoreStyle.Render("[" + string(rune('0'+j.Score/10)) + string(rune('0'+j.Score%10)) + "]")
	companyStr := JobCompanyStyle.Render("@ " + j.Company)
	sourceStr  := JobSourceStyle.Render("(" + j.Source + ")")
	trapStr    := JobTrapsStyle.Render(" [!]")
	traps      := ""
	if j.HasTraps {
		traps = trapStr
	}

	availW := m.width - lipgloss.Width(scoreStr) - lipgloss.Width(companyStr) -
		lipgloss.Width(sourceStr) - lipgloss.Width(traps) - 4
	title := j.Title
	if len(title) > availW && availW > 3 {
		title = title[:availW-3] + "..."
	}
	titleStr := JobItemStyle.Render(title)

	return scoreStr + " " + titleStr + " " + companyStr + " " + sourceStr + traps
}

// ── Status bar — single row ───────────────────────────────────────────────────

func (m Model) renderStatusBar() string {
	keys   := []string{"s", "f", "p", "m", "↑↓", "a", "?", "ctrl+c"}
	labels := []string{"scrape", "filter", "profiles", "emails", "navigate", "apply", "help", "quit"}

	// Footer kbd: same Surface background as the bar — no tint.
	footerKbd := lipgloss.NewStyle().Background(Surface).Foreground(Cyan)
	sp        := lipgloss.NewStyle().Background(Surface).Foreground(Subtle).Render(" ")

	line := ""
	for i, key := range keys {
		if i > 0 {
			line += SepStyle.Render(" │ ")
		}
		line += footerKbd.Render(key) + sp + StatusLabelStyle.Render(labels[i])
	}

	return lipgloss.NewStyle().Background(Surface).Width(m.width).PaddingLeft(2).PaddingRight(2).Render(line)
}