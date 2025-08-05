package views

import (
	"encoding/json"
	"strings"

	"github.com/quinn/gmail-sorter/internal/web/models"
)

// Confirm renders a confirmation prompt for the terminal UI.
func Confirm(actions []models.ActionLink) string {
	var b strings.Builder
	b.WriteString("\n============================\n")
	b.WriteString("      Confirm Action\n")
	b.WriteString("============================\n\n")
	b.WriteString("Are you sure you want to proceed with this action?\n\n")
	if len(actions) > 0 {
		jsonstr, err := json.MarshalIndent(actions[0], "", "  ")
		if err == nil {
			b.WriteString("Action details:\n")
			b.WriteString("----------------------------\n")
			b.WriteString(string(jsonstr))
			b.WriteString("\n----------------------------\n")
		}
	}
	return b.String()
}
