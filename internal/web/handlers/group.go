package handlers

import (
	"net/url"

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
			Path:     "/emails/group-by/domain?val=" + url.QueryEscape(email.FromDomain),
			Label:    "Domain",
			Shortcut: "u",
		},
		{
			Method:   "GET",
			Path:     "/emails/group-by/from?val=" + url.QueryEscape(email.From),
			Label:    "From",
			Shortcut: "f",
		},
		{
			Method:   "GET",
			Path:     "/emails/group-by/to?val=" + url.QueryEscape(email.To),
			Label:    "To",
			Shortcut: "t",
		},
	}

	return pages.GroupEmail(email, actions).Render(c.Request().Context(), c.Response().Writer)
}
