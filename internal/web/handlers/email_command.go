package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

func init() {
	models.Register(EmailCommandAction)
}

var EmailCommandAction models.Action = models.Action{
	ID:               "email-command",
	Method:           "POST",
	Path:             "/emails/:id/command/:command",
	UnwrappedHandler: emailCommand,
	Label:            emailCommandLabel,
}

func emailCommandLabel(link models.ActionLink) string {
	return "Command"
}

// emailCommand handles POST /emails/:id/command/:command
func emailCommand(c echo.Context) error {
	id := c.Param("id")
	api := middleware.GetGmail(c)
	switch c.Param("command") {
	case "skip":
		api.Skip(id)
	case "delete":
		api.Delete(id)
	case "archive":
		api.Archive(id)
	default:
		return fmt.Errorf("invalid command: %s", c.Param("command"))
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
