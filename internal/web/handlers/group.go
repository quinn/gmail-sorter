package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

// GroupEmail handles GET /emails/:id/group
func (h *Handler) GroupEmail(c echo.Context) error {
	id := c.Param("id")

	email, err := h.getEmail(id)
	if err != nil {
		return err
	}

	actions := []models.Action{
		{
			Method:   "GET",
			Path:     "/emails/group-by/domain",
			Label:    "Domain",
			Shortcut: "u",
			Fields: map[string]string{
				"val": email.FromDomain,
			},
		},
		{
			Method:   "GET",
			Path:     "/emails/group-by/from",
			Label:    "From",
			Shortcut: "f",
			Fields: map[string]string{
				"val": email.From,
			},
		},
		{
			Method:   "GET",
			Path:     "/emails/group-by/to",
			Label:    "To",
			Shortcut: "t",
			Fields: map[string]string{
				"val": email.To,
			},
		},
	}

	return pages.GroupEmail(email, actions).Render(c.Request().Context(), c.Response().Writer)
}
