package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

// Page holds the information needed to render a screen.
// It mirrors what the HTML renderer receives.
type Page struct {
	Current models.ActionLink
	Actions []models.ActionLink
	Data    any
}

func RenderPage(link models.ActionLink) (Page, error) {
	page := Page{Current: link}
	ctx := NewContext(&page)
	if err := page.Current.Action().Handler(ctx); err != nil {
		return page, err
	}
	return page, nil
}

// PageMsg is sent to the Bubble Tea program whenever we want to display a new page.
// type PageMsg struct{ Page Page }

// Model implements tea.Model and holds the current page state.
type Model struct {
	page  Page
	err   error
	width int
}

// NewModel constructs a new Bubble Tea model primed with the first page.
func NewModel(initLink models.ActionLink) (Model, error) {
	m := Model{}
	page, err := RenderPage(initLink)
	if err != nil {
		return m, err
	}
	m.page = page
	return m, nil
}

// Init satisfies tea.Model. No initial command yet.
func (m Model) Init() tea.Cmd { return nil }

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	fmt.Printf("Actions (%d): %v\n", len(m.page.Actions), m.page.Actions)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		// Quit on q or ctrl+c for now.
		if msg.String() == "q" || msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		// Handle single-letter shortcut keys for actions
		fmt.Printf("Pressed %s\n", msg.String())
		if len(msg.String()) == 1 && len(m.page.Actions) > 0 {
			shortcuts := models.AssignShortcuts(m.page.Actions)
			pressed := []rune(msg.String())[0]
			for i, key := range shortcuts {
				fmt.Printf("Key %c, Shortcut %c\n", pressed, key)
				if pressed == key {
					fmt.Printf("Pressed %c, triggering action %s\n", pressed, m.page.Actions[i].Action().ID)
					newPage, err := RenderPage(m.page.Actions[i])
					if err != nil {
						m.err = err
						return m, nil
					}
					m.page = newPage
					return m, nil
				}
			}
		}
	}

	return m, nil
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
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}
	if m.width == 0 {
		return "Loading..."
	}
	doc := strings.Builder{}
	// On first render, run the associated handler to populate actions/data.

	// s := fmt.Sprintf("### %s ###\n\n", m.page.Current.Action().ID)
	var actions []string

	shortcuts := models.AssignShortcuts(m.page.Actions)

	viewStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlightColor).
		Padding(0)

	view, err := Render(m.page.Current, m.page.Actions, m.page.Data)
	if err != nil {
		return fmt.Sprintf("### %s failed: %v ###\n\n", m.page.Current.Action().ID, err)
	}

	view = wordwrap.String(view, m.width)

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
	row = lipgloss.JoinVertical(lipgloss.Top, viewStyle.Render(view), row)
	doc.WriteString(row)
	// s += spaceBetween(actions, m.width) + "\n"
	return doc.String()
}
