package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbletea"

	"sprayer/src/api/job"
)

func TestE2E_TUI_EmptyState(t *testing.T) {
	model := NewModel()

	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyNull})
	model = updatedModel.(Model)

	output := model.View()

	if !strings.Contains(output, "No jobs found") {
		t.Errorf("Expected output to contain 'No jobs found', got: %s", output)
	}

	if !strings.Contains(output, "Profile: Default") {
		t.Errorf("Expected output to contain 'Profile: Default', got: %s", output)
	}

	if !strings.Contains(output, "Sprayer") {
		t.Errorf("Expected output to contain 'Sprayer', got: %s", output)
	}
}

func TestE2E_TUI_Navigation(t *testing.T) {
	model := NewModel()
	model.SetJobs([]job.Job{
		{ID: "1", Title: "Software Engineer", Company: "Acme", Score: 85, Source: "remote-ok"},
		{ID: "2", Title: "DevOps Engineer", Company: "Beta", Score: 90, Source: "indeed"},
		{ID: "3", Title: "Backend Developer", Company: "Gamma", Score: 75, Source: "linkedin"},
	})

	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	model = updatedModel.(Model)

	if model.SelectedIndex() != 1 {
		t.Errorf("Expected selectedIndex to be 1 after navigation, got: %d", model.SelectedIndex())
	}
}

func TestE2E_TUI_ViewStateTransitions(t *testing.T) {
	model := NewModel()

	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model = updatedModel.(Model)

	if model.ViewState() != Help {
		t.Errorf("Expected final viewState to be Help, got: %v", model.ViewState())
	}
}

func TestE2E_TUI_Quit(t *testing.T) {
	t.Run("quit with q", func(t *testing.T) {
		model := NewModel()

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

		if cmd == nil {
			t.Errorf("Expected tea.Quit command, got nil")
		}
	})

	t.Run("quit with ctrl+c", func(t *testing.T) {
		model := NewModel()

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

		if cmd == nil {
			t.Errorf("Expected tea.Quit command, got nil")
		}
	})
}

func TestE2E_TUI_JobListRendering(t *testing.T) {
	model := NewModel()
	model.SetJobs([]job.Job{
		{ID: "1", Title: "Software Engineer", Company: "Acme", Score: 85, Source: "remote-ok"},
		{ID: "2", Title: "DevOps Engineer", Company: "Beta", Score: 90, Source: "indeed"},
		{ID: "3", Title: "Backend Developer", Company: "Gamma", Score: 75, Source: "linkedin"},
	})

	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(Model)

	output := model.View()

	if !strings.Contains(output, "Software Engineer") {
		t.Errorf("Expected output to contain 'Software Engineer', got: %s", output)
	}

	if !strings.Contains(output, "DevOps Engineer") {
		t.Errorf("Expected output to contain 'DevOps Engineer', got: %s", output)
	}

	if !strings.Contains(output, "Backend Developer") {
		t.Errorf("Expected output to contain 'Backend Developer', got: %s", output)
	}
}
