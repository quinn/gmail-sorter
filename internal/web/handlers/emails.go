package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

// EmailResponse represents the JSON structure for an email.

// EmailsHandler returns a JSON list of emails from all connected accounts.
func (h *Handler) EmailsHandler(c echo.Context) error {
	// Fetch emails using Gmail API
	var emails []models.EmailResponse
	api := h.spec.GmailService()

	res, err := api.Users.Messages.List("me").MaxResults(50).Do()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	for _, msg := range res.Messages {
		fullMsg, err := api.Users.Messages.Get("me", msg.Id).Format("full").Do()
		if err != nil {
			continue // skip bad messages
		}
		var from, to, subject, date string
		for _, h := range fullMsg.Payload.Headers {
			switch h.Name {
			case "From":
				from = h.Value
			case "To":
				to = h.Value
			case "Subject":
				subject = h.Value
			case "Date":
				date = h.Value
			}
		}
		emails = append(emails, models.EmailResponse{
			ID:       fullMsg.Id,
			ThreadID: fullMsg.ThreadId,
			From:     from,
			To:       to,
			Subject:  subject,
			Date:     date,
			Snippet:  fullMsg.Snippet,
		})
	}

	return pages.Emails(emails).Render(c.Request().Context(), c.Response().Writer)
}
