package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"go.quinn.io/ccf/htmx"
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
	return link.Params[1]
}

// emailCommand handles POST /emails/:id/command/:command
func emailCommand(c echo.Context) error {
	id := c.Param("id")
	api := middleware.GetGmail(c)
	var err error
	switch c.Param("command") {
	case "skip":
		api.Skip(id)
	case "delete":
		err = api.Delete(id)
	case "archive":
		err = api.Archive(id)
	case "open":
		return htmx.Redirect(c, fmt.Sprintf("https://mail.google.com/mail/u/%d/#inbox/%s", api.Account.Index, id))
	case "todo":
		return c.String(http.StatusNotImplemented, "TODO: Move email to Todoist (not yet implemented)")
	default:
		return fmt.Errorf("invalid command: %s", c.Param("command"))
	}

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
