package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

func init() {
	models.Register(EmailSkipAction)
}

var EmailSkipAction models.Action = models.Action{
	Method:           "POST",
	Path:             "/emails/:id/skip",
	Label:            "Skip",
	Shortcut:         "s",
	UnwrappedHandler: skipEmail,
}

// SkipEmail handles POST /emails/:id/skip
func skipEmail(c echo.Context) error {
	id := c.Param("id")
	messages := middleware.GetMessages(c)

	// Remove message with matching id from h.messages
	newMessages := (*messages)[:0]
	for _, m := range *messages {
		if m.Id != id {
			newMessages = append(newMessages, m)
		}
	}
	*messages = newMessages
	return c.Redirect(http.StatusSeeOther, "/")
}
