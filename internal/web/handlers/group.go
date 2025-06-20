package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

// GroupEmailHandler handles GET /emails/:id/group
func (h *Handler) GroupEmailHandler(c echo.Context) error {
	id := c.Param("id")
	return pages.GroupEmail(id).Render(c.Request().Context(), c.Response().Writer)
}
