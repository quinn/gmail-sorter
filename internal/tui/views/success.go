package views

import (
	"encoding/json"
	"strings"

	"github.com/quinn/gmail-sorter/internal/web/models"
)

// Success renders a success message for the terminal UI.
func Success(actions []models.ActionLink, link models.ActionLink) string {
	var b strings.Builder
	b.WriteString("\n============================\n")
	b.WriteString("         Success!\n")
	b.WriteString("============================\n\n")
	b.WriteString("Action completed successfully.\n\n")
	b.WriteString(link.Action().Label(link) + "\n\n")
	jsonstr, err := json.MarshalIndent(link, "", "  ")
	if err == nil {
		b.WriteString("Link details:\n")
		b.WriteString("----------------------------\n")
		b.WriteString(string(jsonstr))
		b.WriteString("\n----------------------------\n")
	}
	return b.String()
}
