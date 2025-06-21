package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// SkipEmailHandler handles POST /emails/:id/skip
func (h *Handler) SkipEmailHandler(c echo.Context) error {
	id := c.Param("id")
	// Remove message with matching id from h.messages
	newMessages := h.messages[:0]
	for _, m := range h.messages {
		if m.Id != id {
			newMessages = append(newMessages, m)
		}
	}
	h.messages = newMessages
	return c.Redirect(http.StatusSeeOther, "/")
}
