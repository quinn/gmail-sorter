package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
	"google.golang.org/api/gmail/v1"
)

func (h *Handler) getEmail(id string) (*gmail.Message, int, error) {
	var msg *gmail.Message
	var idx int
	for i, m := range h.messages {
		if m.Id == id {
			msg = m
			idx = i
			break
		}
	}

	if msg == nil {
		return nil, 0, fmt.Errorf("email not found")
	}

	return msg, idx, nil
}

// EmailHandler renders a single email by ID
func (h *Handler) EmailHandler(c echo.Context) error {
	id := c.Param("id")
	msg, idx, err := h.getEmail(id)
	if err != nil {
		return err
	}

	fullMsg, err := h.spec.GmailService().Users.Messages.Get("me", msg.Id).Format("full").Do()
	if err != nil {
		return fmt.Errorf("failed to get email: %w", err)
	}

	h.messages[idx] = fullMsg
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
