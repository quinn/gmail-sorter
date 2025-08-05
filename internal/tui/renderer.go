package tui

import (
	"fmt"

	"github.com/quinn/gmail-sorter/internal/tui/views"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

func Render(current models.ActionLink, actions []models.ActionLink, data interface{}) (string, error) {
	// Render based on actionID
	switch current.Action().ID {
	case "confirm":
		return views.Confirm(actions), nil
	case "success":
		link, ok := data.(models.ActionLink)
		if !ok {
			return "", fmt.Errorf("expected ActionLink, got %T", data)
		}
		return views.Success(actions, link), nil
	}

	// Handle data-based rendering
	switch typed := data.(type) {
	case models.FiltersPageData:
		return views.Filters(typed.AccountID, typed.Filters, actions), nil
	case models.GroupByPageData:
		return views.GroupBy(typed.GroupType, typed.Value, typed.Emails, actions), nil
	case models.GroupByDeleteSuccessPageData:
		return views.GroupByDeleteSuccess(actions, typed.GroupType, typed.Value, typed.Count), nil
	case models.EmailResponse:
		return views.Email(typed, actions), nil
	}

	return "", fmt.Errorf("no TUI renderer implemented for action ID: %s", current.Action().ID)

}
