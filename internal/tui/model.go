package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

// Page holds the information needed to render a screen.
// It mirrors what the HTML renderer receives.
type Page struct {
	Current models.ActionLink
	Actions []models.ActionLink
	Data    any
}

// PageMsg is sent to the Bubble Tea program whenever we want to display a new page.
type PageMsg struct{ Page Page }

// Model implements tea.Model and holds the current page state.
type Model struct {
	page  Page
	width int
}

// NewModel constructs a new Bubble Tea model primed with the first page.
func NewModel(p Page) Model { return Model{page: p} }

// Init satisfies tea.Model. No initial command yet.
func (m Model) Init() tea.Cmd { return nil }

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		// Quit on q or ctrl+c for now.
		if msg.String() == "q" || msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case PageMsg:
		m.page = msg.Page
	}

	return m, nil
}

func spaceBetween(items []string, totalWidth int) string {
	n := len(items)
	if n == 0 {
		return ""
	}
	sum := 0
	for _, it := range items {
		sum += lipgloss.Width(it)
	}
	if n == 1 || totalWidth <= sum {
		return strings.Join(items, " ")
	}
	remaining := totalWidth - sum
	gapCount := n - 1
	base := remaining / gapCount
	extra := remaining % gapCount

	var b strings.Builder
	for i, it := range items {
		b.WriteString(it)
		if i < gapCount {
			gap := base
			if i < extra {
				gap++ // distribute leftovers to the left
			}
			b.WriteString(strings.Repeat(" ", gap))
		}
	}
	return b.String()
}

// View renders the current page to a string.
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}
	// On first render, run the associated handler to populate actions/data.
	ctx := NewContext(&m)
	if err := m.page.Current.Action().Handler(ctx); err != nil {
		return fmt.Sprintf("### %s failed: %v ###\n\n", m.page.Current.Action().ID, err)
	}

	s := fmt.Sprintf("### %s ###\n\n", m.page.Current.Action().ID)
	var actions []string
	if len(m.page.Actions) > 0 {
		s += "Available actions:\n"
		for i, a := range m.page.Actions {
			actions = append(actions, fmt.Sprintf("%d. %s", i+1, a.Action().Label(a)))
		}
	}
	s += spaceBetween(actions, m.width) + "\n"
	s += "\nPress q to quit."
	return s
}
