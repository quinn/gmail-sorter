package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// NewProgram constructs a Bubble Tea program starting with the given initial Page.
// The caller should run the program with prog.Run().
func NewProgram(start Page) *tea.Program {
	m := NewModel(start)
	return tea.NewProgram(m, tea.WithAltScreen())
}
