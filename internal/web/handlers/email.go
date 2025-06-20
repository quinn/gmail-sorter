package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
	"google.golang.org/api/gmail/v1"
)

// emailNavIDs returns the previous and next email IDs given a list and a current ID.
func emailNavIDs(emails []*gmail.Message, currentID string) (prev, next string) {
	for i, msg := range emails {
		if msg.Id == currentID {
			if i > 0 {
				prev = emails[i-1].Id
			}
			if i < len(emails)-1 {
				next = emails[i+1].Id
			}
			break
		}
	}
	return
}

// EmailHandler renders a single email by ID, with prev/next navigation
func (h *Handler) EmailHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusBadRequest, "Missing email ID")
	}
	api := h.spec.GmailService()

	// Fetch the list of messages for navigation
	listRes, err := api.Users.Messages.List("me").MaxResults(50).Do()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to list emails")
	}
	prevID, nextID := emailNavIDs(listRes.Messages, id)

	msg, err := api.Users.Messages.Get("me", id).Format("full").Do()
	if err != nil {
		return c.String(http.StatusNotFound, "Email not found")
	}
	var from, to, subject, date string
	for _, h := range msg.Payload.Headers {
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
	email := models.EmailResponse{
		ID:       msg.Id,
		ThreadID: msg.ThreadId,
		From:     from,
		To:       to,
		Subject:  subject,
		Date:     date,
		Snippet:  msg.Snippet,
		PrevID:   prevID,
		NextID:   nextID,
	}
	return pages.Email(email).Render(c.Request().Context(), c.Response().Writer)
}
