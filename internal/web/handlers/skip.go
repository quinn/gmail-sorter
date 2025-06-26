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
	ID:               "email-skip",
	Method:           "POST",
	Path:             "/emails/:id/skip",
	UnwrappedHandler: skipEmail,
	Label:            skipLabel,
}

func skipLabel(link models.ActionLink) string {
	return "Skip"
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
