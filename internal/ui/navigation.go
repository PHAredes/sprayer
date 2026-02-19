package ui

import "strings"

// NavigationHistory manages state history for proper back navigation
type NavigationHistory struct {
	history []AppState
	current int
}

// NewNavigationHistory creates a new navigation history
func NewNavigationHistory() *NavigationHistory {
	return &NavigationHistory{
		history: []AppState{StateList}, // Start with main list
		current: 0,
	}
}

// Push adds a new state to history
func (nh *NavigationHistory) Push(state AppState) {
	// Don't push the same state consecutively
	if len(nh.history) > 0 && nh.history[nh.current] == state {
		return
	}

	// Remove states after current position (when user navigated back)
	nh.history = nh.history[:nh.current+1]

	// Add new state
	nh.history = append(nh.history, state)
	nh.current++

	// Limit history to prevent memory issues
	if len(nh.history) > 20 {
		nh.history = nh.history[len(nh.history)-20:]
		nh.current = 19
	}
}

// Back returns to the previous state
func (nh *NavigationHistory) Back() AppState {
	if nh.current > 0 {
		nh.current--
		return nh.history[nh.current]
	}
	return nh.history[0] // Stay at root if can't go back
}

// CanBack returns whether we can go back
func (nh *NavigationHistory) CanBack() bool {
	return nh.current > 0
}

// Current returns the current state
func (nh *NavigationHistory) Current() AppState {
	if nh.current >= 0 && nh.current < len(nh.history) {
		return nh.history[nh.current]
	}
	return StateList
}

// Clear resets navigation history
func (nh *NavigationHistory) Clear() {
	nh.history = []AppState{StateList}
	nh.current = 0
}

// GetBreadcrumb returns navigation breadcrumb for context
func (nh *NavigationHistory) GetBreadcrumb() string {
	if len(nh.history) <= 1 {
		return ""
	}

	var breadcrumb []string
	for i := 0; i <= nh.current && i < len(nh.history); i++ {
		state := nh.history[i]
		name := getStateName(state)
		if name != "" {
			breadcrumb = append(breadcrumb, name)
		}
	}

	if len(breadcrumb) <= 1 {
		return ""
	}

	return strings.Join(breadcrumb[len(breadcrumb)-2:], " → ")
}

// getStateName returns a human-readable name for a state
func getStateName(state AppState) string {
	switch state {
	case StateList:
		return "Jobs"
	case StateDetail:
		return "Details"
	case StateFilters:
		return "Filters"
	case StateProfiles:
		return "Profiles"
	case StateReview:
		return "Review"
	case StateHelp:
		return "Help"
	case StateScraping:
		return "Scraping"
	default:
		return ""
	}
}
