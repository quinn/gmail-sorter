package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/quinn/gmail-sorter/internal/web/models"
)

func init() {
	models.Register(ConfirmAction)
}

var ConfirmAction models.Action = models.Action{
	ID:      "confirm",
	Method:  "GET",
	Path:    "/confirm",
	Handler: confirm,
	Label:   confirmLabel,
}

func confirmLabel(link models.ActionLink) string {
	return "Confirm"
}

// Confirm handles the /confirm endpoint
func confirm(c models.Context) error {
	link := c.QueryParam("link")
	var actionObj models.ActionLink
	if err := json.Unmarshal([]byte(link), &actionObj); err != nil {
		return fmt.Errorf("failed to unmarshal action: %w", err)
	}

	actionObj.Confirm = false

	return c.Render(nil, nil)
}
