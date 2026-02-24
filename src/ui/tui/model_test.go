package tui

import (
	"testing"

	"github.com/charmbracelet/bubbletea"

	"sprayer/src/api/job"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if len(m.jobs) != 0 {
		t.Errorf("expected jobs to be empty slice, got %v", m.jobs)
	}

	if m.selectedIndex != 0 {
		t.Errorf("expected selectedIndex to be 0, got %d", m.selectedIndex)
	}

	if m.profileName != "Default" {
		t.Errorf("expected profileName to be 'Default', got %s", m.profileName)
	}

	if m.viewState != EmptyState {
		t.Errorf("expected viewState to be EmptyState, got %v", m.viewState)
	}

	if m.width != 80 {
		t.Errorf("expected width to be 80, got %d", m.width)
	}

	if m.height != 24 {
		t.Errorf("expected height to be 24, got %d", m.height)
	}
}

func TestModel_Init(t *testing.T) {
	m := NewModel()
	cmd := m.Init()

	if cmd != nil {
		t.Errorf("expected Init to return nil, got %v", cmd)
	}
}

func TestModel_Update_Navigation(t *testing.T) {
	tests := []struct {
		name            string
		jobs            []job.Job
		key             string
		expectIndex     int
		expectViewState ViewState
	}{
		{
			name:            "j increments index when jobs exist",
			jobs:            []job.Job{{ID: "1"}, {ID: "2"}, {ID: "3"}},
			key:             "j",
			expectIndex:     1,
			expectViewState: JobList,
		},
		{
			name:            "down arrow increments index when jobs exist",
			jobs:            []job.Job{{ID: "1"}, {ID: "2"}, {ID: "3"}},
			key:             "↓",
			expectIndex:     1,
			expectViewState: JobList,
		},
		{
			name:            "k decrements index when jobs exist and index > 0",
			jobs:            []job.Job{{ID: "1"}, {ID: "2"}, {ID: "3"}},
			key:             "k",
			expectIndex:     0,
			expectViewState: JobList,
		},
		{
			name:            "up arrow decrements index when jobs exist and index > 0",
			jobs:            []job.Job{{ID: "1"}, {ID: "2"}, {ID: "3"}},
			key:             "↑",
			expectIndex:     0,
			expectViewState: JobList,
		},
		{
			name:            "j does nothing when jobs is empty",
			jobs:            []job.Job{},
			key:             "j",
			expectIndex:     0,
			expectViewState: EmptyState,
		},
		{
			name:            "down arrow does nothing when jobs is empty",
			jobs:            []job.Job{},
			key:             "↓",
			expectIndex:     0,
			expectViewState: EmptyState,
		},
		{
			name:            "k does nothing when jobs is empty",
			jobs:            []job.Job{},
			key:             "k",
			expectIndex:     0,
			expectViewState: EmptyState,
		},
		{
			name:            "up arrow does nothing when jobs is empty",
			jobs:            []job.Job{},
			key:             "↑",
			expectIndex:     0,
			expectViewState: EmptyState,
		},
		{
			name:            "k does nothing when at index 0",
			jobs:            []job.Job{{ID: "1"}, {ID: "2"}, {ID: "3"}},
			key:             "k",
			expectIndex:     0,
			expectViewState: JobList,
		},
		{
			name:            "up arrow does nothing when at index 0",
			jobs:            []job.Job{{ID: "1"}, {ID: "2"}, {ID: "3"}},
			key:             "↑",
			expectIndex:     0,
			expectViewState: JobList,
		},
		{
			name:            "j stays within bounds at last index",
			jobs:            []job.Job{{ID: "1"}},
			key:             "j",
			expectIndex:     0,
			expectViewState: JobList,
		},
		{
			name:            "down arrow stays within bounds at last index",
			jobs:            []job.Job{{ID: "1"}},
			key:             "↓",
			expectIndex:     0,
			expectViewState: JobList,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()
			m.jobs = tt.jobs

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			newModel, _ := m.Update(msg)
			model := newModel.(Model)

			if model.selectedIndex != tt.expectIndex {
				t.Errorf("expected selectedIndex %d, got %d", tt.expectIndex, model.selectedIndex)
			}

			if model.viewState != tt.expectViewState {
				t.Errorf("expected viewState %v, got %v", tt.expectViewState, model.viewState)
			}
		})
	}
}

func TestModel_Update_ViewStates(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		expectState ViewState
	}{
		{
			name:        "s changes to Scraping",
			key:         "s",
			expectState: Scraping,
		},
		{
			name:        "f changes to Filter",
			key:         "f",
			expectState: Filter,
		},
		{
			name:        "p changes to Profiles",
			key:         "p",
			expectState: Profiles,
		},
		{
			name:        "m changes to Emails",
			key:         "m",
			expectState: Emails,
		},
		{
			name:        "? changes to Help",
			key:         "?",
			expectState: Help,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			newModel, _ := m.Update(msg)
			model := newModel.(Model)

			if model.viewState != tt.expectState {
				t.Errorf("expected viewState %v, got %v", tt.expectState, model.viewState)
			}
		})
	}
}

func TestModel_Update_Quit(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "ctrl+c returns quit",
			key:  "ctrl+c",
		},
		{
			name: "q returns quit",
			key:  "q",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			_, cmd := m.Update(msg)

			if cmd == nil {
				t.Error("expected non-nil command")
				return
			}
			if _, ok := cmd().(tea.QuitMsg); !ok {
				t.Errorf("expected tea.QuitMsg, got %T", cmd())
			}
		})
	}
}

func TestModel_Update_WindowSize(t *testing.T) {
	m := NewModel()

	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	newModel, _ := m.Update(msg)
	model := newModel.(Model)

	if model.width != 120 {
		t.Errorf("expected width 120, got %d", model.width)
	}

	if model.height != 40 {
		t.Errorf("expected height 40, got %d", model.height)
	}
}

func TestModel_View_EmptyState(t *testing.T) {
	m := NewModel()
	m.width = 80
	m.height = 24

	view := m.View()

	if len(view) == 0 {
		t.Error("expected View() to return non-empty string")
	}

	if !contains(view, "No jobs found") {
		t.Error("expected View() to contain 'No jobs found'")
	}
}

func TestModel_View_JobList(t *testing.T) {
	m := NewModel()
	m.jobs = []job.Job{
		{
			ID:       "1",
			Title:    "Software Engineer",
			Company:  "TechCorp",
			Source:   "Indeed",
			Score:    85,
			HasTraps: false,
		},
		{
			ID:       "2",
			Title:    "Backend Developer",
			Company:  "StartupXYZ",
			Source:   "LinkedIn",
			Score:    72,
			HasTraps: true,
		},
	}
	m.width = 80
	m.height = 24

	view := m.View()

	if !contains(view, "Software Engineer") {
		t.Error("expected View() to contain job title 'Software Engineer'")
	}

	if !contains(view, "TechCorp") {
		t.Error("expected View() to contain company 'TechCorp'")
	}

	if !contains(view, "Indeed") {
		t.Error("expected View() to contain source 'Indeed'")
	}

	if !contains(view, "StartupXYZ") {
		t.Error("expected View() to contain company 'StartupXYZ'")
	}

	if !contains(view, "LinkedIn") {
		t.Error("expected View() to contain source 'LinkedIn'")
	}
}

func TestModel_View_TopBar(t *testing.T) {
	m := NewModel()
	m.profileName = "TestProfile"
	m.jobs = []job.Job{{ID: "1"}, {ID: "2"}}
	m.width = 80
	m.height = 24

	view := m.View()

	if !contains(view, "TestProfile") {
		t.Error("expected View() to contain profile name 'TestProfile'")
	}

	if !contains(view, "Sprayer") {
		t.Error("expected View() to contain 'Sprayer'")
	}

	if !contains(view, "Jobs:") {
		t.Error("expected View() to contain 'Jobs:'")
	}
}

func TestModel_View_StatusBar(t *testing.T) {
	m := NewModel()
	m.width = 80
	m.height = 24

	view := m.View()

	if !contains(view, "s") {
		t.Error("expected View() to contain key hint 's'")
	}

	if !contains(view, "f") {
		t.Error("expected View() to contain key hint 'f'")
	}

	if !contains(view, "p") {
		t.Error("expected View() to contain key hint 'p'")
	}

	if !contains(view, "m") {
		t.Error("expected View() to contain key hint 'm'")
	}

	if !contains(view, "?") {
		t.Error("expected View() to contain key hint '?'")
	}

	if !contains(view, "q") {
		t.Error("expected View() to contain key hint 'q'")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
