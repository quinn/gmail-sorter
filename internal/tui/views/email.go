package views

import (
	"fmt"
	"strings"

	"github.com/quinn/gmail-sorter/internal/web/models"
)

// Email renders an email page for the terminal UI.
func Email(email models.EmailResponse, actions []models.ActionLink) string {
	var b strings.Builder
	b.WriteString("\n============================\n")
	b.WriteString("            Email\n")
	b.WriteString("============================\n\n")
	b.WriteString(fmt.Sprintf("From:    %s\n", strings.Join(email.From, ", ")))
	b.WriteString(fmt.Sprintf("To:      %s\n", strings.Join(email.To, ", ")))
	b.WriteString(fmt.Sprintf("Subject: %s\n", email.Subject))
	b.WriteString(fmt.Sprintf("Date:    %s\n", email.Date))
	b.WriteString("\n----------------------------\n")
	b.WriteString(email.Snippet)
	b.WriteString("\n----------------------------\n")
	return b.String()
}
