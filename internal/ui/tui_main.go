package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// TUI represents the main terminal user interface
type TUI struct {
	model *Model
}

// NewTUI creates a new terminal user interface
func NewTUI() (*TUI, error) {
	model, err := NewModel()
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %w", err)
	}

	return &TUI{
		model: &model,
	}, nil
}

// Run starts the TUI application
func (t *TUI) Run() error {
	p := tea.NewProgram(
		t.model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithInput(os.Stdin),
		tea.WithOutput(os.Stdout),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}

// InitializeTUI creates and runs the TUI
func InitializeTUI() error {
	tui, err := NewTUI()
	if err != nil {
		return fmt.Errorf("failed to create TUI: %w", err)
	}

	return tui.Run()
}
