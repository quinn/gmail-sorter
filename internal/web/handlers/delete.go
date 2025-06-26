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
	ID:               "email-delete",
	Method:           "POST",
	Path:             "/emails/:id/delete",
	UnwrappedHandler: deleteEmail,
	Label:            deleteLabel,
}

func deleteLabel(link models.ActionLink) string {
	return "Delete"
}

// DeleteEmail handles POST /emails/:id/delete
func deleteEmail(c echo.Context) error {
	id := c.Param("id")
	// TODO: Implement business logic
	return c.String(http.StatusOK, "Delete action for email "+id)
}
