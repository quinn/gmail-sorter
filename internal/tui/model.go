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

var (
	bottomJoint    = "┴"
	topJoint       = "┬"
	highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	buttonStyle = lipgloss.
			NewStyle().
			Border(lipgloss.RoundedBorder(), true, false, true, true).
			BorderForeground(highlightColor).
			Padding(0)
)

// View renders the current page to a string.
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}
	doc := strings.Builder{}
	// On first render, run the associated handler to populate actions/data.
	ctx := NewContext(&m)
	if err := m.page.Current.Action().Handler(ctx); err != nil {
		return fmt.Sprintf("### %s failed: %v ###\n\n", m.page.Current.Action().ID, err)
	}

	// s := fmt.Sprintf("### %s ###\n\n", m.page.Current.Action().ID)
	var actions []string

	shortcuts := models.AssignShortcuts(m.page.Actions)

	for i, act := range m.page.Actions {
		isFirst := i == 0
		isLast := i == len(m.page.Actions)-1
		// fmt.Sprintf("%d %t %t\n", i, isFirst, isLast)
		style := buttonStyle
		text := fmt.Sprintf("%s [%s]", act.Action().Label(act), string(shortcuts[i]))

		border, _, _, _, _ := style.GetBorder()
		if !isFirst {
			border.BottomLeft = bottomJoint
			border.TopLeft = topJoint
		}
		style = style.Border(border)
		style = style.BorderRight(false)

		if isLast {
			style = style.BorderRight(true)
		}
		style = style.Width(m.width/len(m.page.Actions) - style.GetHorizontalFrameSize()).Align(lipgloss.Center)
		actions = append(actions, style.Render(text))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, actions...)
	doc.WriteString(row)
	// s += spaceBetween(actions, m.width) + "\n"
	return doc.String()
}
