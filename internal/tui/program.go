package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

// NewProgram constructs a Bubble Tea program starting with the given initial Page.
// The caller should run the program with prog.Run().
func NewProgram(start models.ActionLink) (*tea.Program, error) {
	m, err := NewModel(start)
	if err != nil {
		return nil, err
	}
	return tea.NewProgram(m, tea.WithAltScreen()), nil
}
