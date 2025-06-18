package handlers

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

// EmailResponse represents the JSON structure for an email.
type EmailResponse struct {
	ID      string `json:"id"`
	ThreadID string `json:"threadId"`
	From    string `json:"from"`
	To      string `json:"to,omitempty"`
	Subject string `json:"subject"`
	Date    string `json:"date"`
	Snippet string `json:"snippet"`
}

// EmailsHandler returns a JSON list of emails from all connected accounts.
func (h *Handler) EmailsHandler(c echo.Context) error {
	if h.spec == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "spec not set in handler"})
	}

	// Fetch emails using Gmail API
	var emails []EmailResponse
	pageToken := ""
	api := h.spec.GmailService()
	for {
		res, err := api.Users.Messages.List("me").MaxResults(50).PageToken(pageToken).Do()
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
			emails = append(emails, EmailResponse{
				ID:      fullMsg.Id,
				ThreadID: fullMsg.ThreadId,
				From:    from,
				To:      to,
				Subject: subject,
				Date:    date,
				Snippet: fullMsg.Snippet,
			})
		}
		if res.NextPageToken == "" {
			break
		}
		pageToken = res.NextPageToken
	}

	return c.JSON(http.StatusOK, emails)
}


