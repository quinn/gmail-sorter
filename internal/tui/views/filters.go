package views

import (
	"fmt"
	"strings"

	"github.com/quinn/gmail-sorter/internal/web/models"
	"google.golang.org/api/gmail/v1"
)

// Filters renders the filters page for the terminal UI.
func Filters(accountID uint, filters []*gmail.Filter, actions []models.ActionLink) string {
	var b strings.Builder
	b.WriteString("\n============================\n")
	b.WriteString("          Filters\n")
	b.WriteString("============================\n\n")
	b.WriteString("ID           | Criteria                  | Actions\n")
	b.WriteString("-------------+---------------------------+-----------------------\n")
	for _, filter := range filters {
		criteria := describeCriteria(filter.Criteria)
		actionsDesc := describeAction(filter.Action)
		b.WriteString(fmt.Sprintf("%-12s | %-25s | %s\n", filter.Id, criteria, actionsDesc))
	}
	return b.String()
}

// describeCriteria provides a simple string summary of filter criteria for TUI.
func describeCriteria(c *gmail.FilterCriteria) string {
	if c == nil {
		return "-"
	}
	var parts []string
	if c.From != "" {
		parts = append(parts, "From: "+c.From)
	}
	if c.To != "" {
		parts = append(parts, "To: "+c.To)
	}
	if c.Subject != "" {
		parts = append(parts, "Subj: "+c.Subject)
	}
	if c.Query != "" {
		parts = append(parts, "Query: "+c.Query)
	}
	return strings.Join(parts, ", ")
}

// describeAction provides a simple string summary of filter action for TUI.
func describeAction(a *gmail.FilterAction) string {
	if a == nil {
		return "-"
	}
	var parts []string
	if a.AddLabelIds != nil && len(a.AddLabelIds) > 0 {
		parts = append(parts, "Add: "+strings.Join(a.AddLabelIds, ","))
	}
	if a.RemoveLabelIds != nil && len(a.RemoveLabelIds) > 0 {
		parts = append(parts, "Remove: "+strings.Join(a.RemoveLabelIds, ","))
	}
	if a.Forward != "" {
		parts = append(parts, "Fwd: "+a.Forward)
	}
	return strings.Join(parts, ", ")
}
