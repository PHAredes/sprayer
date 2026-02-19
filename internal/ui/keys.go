package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up            key.Binding
	Down          key.Binding
	Left          key.Binding
	Right         key.Binding
	Enter         key.Binding
	Esc           key.Binding
	Quit          key.Binding
	Scrape        key.Binding
	Filter        key.Binding
	Profiles      key.Binding
	Apply         key.Binding
	Send          key.Binding
	Help          key.Binding
	Tab           key.Binding
	Back          key.Binding
	Sort          key.Binding
	ClearFilter   key.Binding
	NewProfile    key.Binding
	EditProfile   key.Binding
	DeleteProfile key.Binding
	ImportProfile key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns all keybindings for the expanded help view.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Esc, k.Tab, k.Back},
		{k.Scrape, k.Filter, k.Profiles, k.Apply},
		{k.Sort, k.ClearFilter, k.Send, k.Help},
		{k.Quit},
	}
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
	Scrape: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "scrape"),
	),
	Filter: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "filter"),
	),
	Profiles: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "profiles"),
	),
	Apply: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "apply"),
	),
	Send: key.NewBinding(
		key.WithKeys("S"),
		key.WithHelp("S", "send (smtp)"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
	Back: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev field"),
	),
	Sort: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "sort"),
	),
	ClearFilter: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "clear filter"),
	),
	NewProfile: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new profile"),
	),
	EditProfile: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit profile"),
	),
	DeleteProfile: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete profile"),
	),
	ImportProfile: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "import profile"),
	),
}
