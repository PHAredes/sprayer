package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// StatusBar represents a stable status bar component with fixed positioning
type StatusBar struct {
	mu           sync.RWMutex
	message      string
	timer        int
	lastTime     string
	timeUpdated  time.Time
	width        int
	cachedRender string
	needsRefresh bool
}

// NewStatusBar creates a new stable status bar
func NewStatusBar() *StatusBar {
	return &StatusBar{
		message:      "Ready",
		lastTime:     time.Now().Format("15:04"),
		timeUpdated:  time.Now(),
		needsRefresh: true,
	}
}

// StatusBarInfo contains information to display in the status bar
type StatusBarInfo struct {
	Message     string
	Profile     string
	JobCount    int
	TotalJobs   int
	StatusTimer int
}

// UpdateSize updates the status bar width and triggers refresh
func (s *StatusBar) UpdateSize(width int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.width != width {
		s.width = width
		s.needsRefresh = true
	}
}

// View renders the status bar with stable formatting
func (s *StatusBar) View(width int, info StatusBarInfo) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update message if provided
	if info.Message != "" {
		s.message = info.Message
		s.timer = info.StatusTimer
		s.needsRefresh = true
	} else if s.timer > 0 {
		s.timer--
		s.needsRefresh = true
	} else if s.timer == 0 && s.message != "Ready" {
		s.message = "Ready"
		s.needsRefresh = true
	}

	// Update time only every 30 seconds to prevent constant re-renders
	if time.Since(s.timeUpdated) > 30*time.Second {
		s.lastTime = time.Now().Format("15:04")
		s.timeUpdated = time.Now()
		s.needsRefresh = true
	}

	// Use cached render if nothing changed and width is same
	if !s.needsRefresh && s.width == width && s.cachedRender != "" {
		return s.cachedRender
	}

	// Build status bar content with fixed widths
	left := s.buildLeftSection(info)
	center := s.buildCenterSection()
	right := s.buildRightSection()

	// Calculate fixed widths to prevent jumping
	leftWidth := 25  // Fixed width for profile/jobs info
	rightWidth := 20 // Fixed width for time/help
	centerWidth := width - leftWidth - rightWidth - 4

	// Ensure minimum widths
	if centerWidth < 10 {
		centerWidth = 10
		leftWidth = (width - 30) / 2
		rightWidth = (width - 30) / 2
	}

	// Style each section with fixed widths
	leftStyled := Styles.StatusText.Width(leftWidth).Render(left)
	centerStyled := Styles.StatusText.Width(centerWidth).Align(lipgloss.Center).Render(center)
	rightStyled := Styles.StatusText.Width(rightWidth).Align(lipgloss.Right).Render(right)

	// Cache the render
	s.cachedRender = Styles.StatusBar.Width(width).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, leftStyled, centerStyled, rightStyled),
	)
	s.needsRefresh = false
	s.width = width

	return s.cachedRender
}

func (s *StatusBar) buildLeftSection(info StatusBarInfo) string {
	parts := []string{}

	// Profile name (truncated if too long)
	if info.Profile != "" {
		profile := info.Profile
		if len(profile) > 12 {
			profile = profile[:9] + "..."
		}
		parts = append(parts, profile)
	}

	// Job counts
	if info.JobCount > 0 {
		parts = append(parts, fmt.Sprintf("%d/%d", info.JobCount, info.TotalJobs))
	}

	if len(parts) == 0 {
		return "Ready"
	}

	return strings.Join(parts, " • ")
}

func (s *StatusBar) buildCenterSection() string {
	// Truncate message to fit fixed width
	msg := s.message
	maxLen := 30 // Fixed max length for center section
	if len(msg) > maxLen {
		msg = msg[:maxLen-3] + "..."
	}
	return msg
}

func (s *StatusBar) buildRightSection() string {
	return fmt.Sprintf("%s • ? help", s.lastTime)
}

// SetMessage updates the status message and triggers refresh
func (s *StatusBar) SetMessage(message string, duration int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
	s.timer = duration
	s.needsRefresh = true
}

// SetPermanentMessage sets a message that doesn't timeout
func (s *StatusBar) SetPermanentMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
	s.timer = -1
	s.needsRefresh = true
}

// ClearMessage clears the current message
func (s *StatusBar) ClearMessage() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = "Ready"
	s.timer = 0
	s.needsRefresh = true
}

// ForceRefresh forces a complete refresh on next render
func (s *StatusBar) ForceRefresh() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.needsRefresh = true
	s.cachedRender = ""
}
