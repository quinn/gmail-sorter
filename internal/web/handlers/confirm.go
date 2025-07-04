package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

func init() {
	models.Register(ConfirmAction)
}

var ConfirmAction models.Action = models.Action{
	ID:               "confirm",
	Method:           "GET",
	Path:             "/confirm",
	UnwrappedHandler: confirm,
	Label:            confirmLabel,
}

func confirmLabel(link models.ActionLink) string {
	return "Confirm"
}

// Confirm handles the /confirm endpoint
func confirm(c echo.Context) error {
	link := c.QueryParam("link")
	var actionObj models.ActionLink
	if err := json.Unmarshal([]byte(link), &actionObj); err != nil {
		return fmt.Errorf("failed to unmarshal action: %w", err)
	}

	actionObj.Confirm = false

	return pages.Confirm([]models.ActionLink{actionObj}).Render(c.Request().Context(), c.Response().Writer)
}
