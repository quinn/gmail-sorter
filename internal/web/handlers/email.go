package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
	"google.golang.org/api/gmail/v1"
)

func (h *Handler) getEmail(id string) (email models.EmailResponse, err error) {
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
		return email, fmt.Errorf("email not found")
	}

	if msg.Payload == nil {
		fullMsg, err := h.spec.GmailService().Users.Messages.Get("me", msg.Id).Format("full").Do()
		if err != nil {
			return email, fmt.Errorf("failed to get email: %w", err)
		}
		msg = fullMsg
		h.messages[idx] = fullMsg
	}

	email = models.FromGmailMessage(msg)
	return email, nil
}

// Email renders a single email by ID
func (h *Handler) Email(c echo.Context) error {
	id := c.Param("id")
	email, err := h.getEmail(id)
	if err != nil {
		return err
	}

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
