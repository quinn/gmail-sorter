package views

import (
	"fmt"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

// GroupByDeleteSuccess renders a success message for group-by delete actions in the terminal UI.
func GroupByDeleteSuccess(actions []models.ActionLink, groupType, val string, count int) string {
	return fmt.Sprintf(`
============================
 Deleted %d emails grouped by %s: %s
============================
The selected emails have been moved to Trash.
`, count, groupType, val)
}
