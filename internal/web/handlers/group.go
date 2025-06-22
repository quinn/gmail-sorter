package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

// GroupEmail handles GET /emails/:id/group
func (h *Handler) GroupEmail(c echo.Context) error {
	id := c.Param("id")

	msg, _, err := h.getEmail(id)
	if err != nil {
		return err
	}

	email := models.FromGmailMessage(msg)

	actions := []models.Action{
		{
			Method:   "GET",
			Path:     "/emails/" + id + "/group/by/domain",
			Label:    "Domain",
			Shortcut: "u",
		},
		{
			Method:   "GET",
			Path:     "/emails/" + id + "/group/by/from",
			Label:    "From",
			Shortcut: "f",
		},
		{
			Method:   "GET",
			Path:     "/emails/" + id + "/group/by/to",
			Label:    "To",
			Shortcut: "t",
		},
	}

	return pages.GroupEmail(email, actions).Render(c.Request().Context(), c.Response().Writer)
}
