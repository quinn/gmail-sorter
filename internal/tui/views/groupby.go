package views

import (
	"fmt"
	"strings"

	"github.com/quinn/gmail-sorter/internal/web/models"
)

// GroupBy renders the grouped email page for the terminal UI.
func GroupBy(groupType, val string, emails []models.EmailResponse, actions []models.ActionLink) string {
	var b strings.Builder
	b.WriteString("\n============================\n")
	b.WriteString(fmt.Sprintf(" Grouped by %s\n", groupType))
	b.WriteString("============================\n")
	b.WriteString(fmt.Sprintf("Group value: %s\n\n", val))
	b.WriteString("From            | To              | Subject         | Date       | Snippet\n")
	b.WriteString("----------------+-----------------+-----------------+------------+---------------------\n")
	for _, email := range emails {
		from := strings.Join(email.From, ", ")
		to := strings.Join(email.To, ", ")
		subj := email.Subject
		date := email.Date
		snippet := email.Snippet
		b.WriteString(fmt.Sprintf("%-15s | %-15s | %-15s | %-10s | %.20s\n", from, to, subj, date, snippet))
	}
	return b.String()
}
