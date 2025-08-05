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

	// For now, only confirm and success are supported for TUI.
	return "", fmt.Errorf("no TUI renderer implemented for action ID: %s", current.Action().ID)

}
