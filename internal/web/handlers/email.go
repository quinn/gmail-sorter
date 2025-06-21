package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
	"google.golang.org/api/gmail/v1"
)

// EmailHandler renders a single email by ID
func (h *Handler) EmailHandler(c echo.Context) error {
	id := c.Param("id")
	var msg *gmail.Message
	for _, m := range h.messages {
		if m.Id == id {
			msg = m
			break
		}
	}
	if msg == nil {
		return c.String(http.StatusNotFound, "Email not found")
	}

	fullMsg, err := h.spec.GmailService().Users.Messages.Get("me", msg.Id).Format("full").Do()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get email: "+err.Error())
	}

	email := models.FromGmailMessage(fullMsg)

	actions := []models.Action{
		{
			Method:   "GET",
			Path:     "/emails/" + id + "/group",
			Label:    "Group",
			Shortcut: "g",
		},
		{
			Method:   "POST",
			Path:     "/emails/" + id + "/skip",
			Label:    "Skip",
			Shortcut: "s",
		},
		{
			Method:   "POST",
			Path:     "/emails/" + id + "/delete",
			Label:    "Delete",
			Shortcut: "d",
		},
	}
	return pages.Email(email, actions).Render(c.Request().Context(), c.Response().Writer)
}
