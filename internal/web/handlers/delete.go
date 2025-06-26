package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

func init() {
	models.Register(EmailDeleteAction)
}

var EmailDeleteAction models.Action = models.Action{
	Method:           "POST",
	Path:             "/emails/:id/delete",
	Label:            "Delete",
	Shortcut:         "d",
	UnwrappedHandler: deleteEmail,
}

// DeleteEmail handles POST /emails/:id/delete
func deleteEmail(c echo.Context) error {
	id := c.Param("id")
	// TODO: Implement business logic
	return c.String(http.StatusOK, "Delete action for email "+id)
}
