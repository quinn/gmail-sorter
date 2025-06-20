package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// SkipEmailHandler handles POST /emails/:id/skip
func (h *Handler) SkipEmailHandler(c echo.Context) error {
	id := c.Param("id")
	// TODO: Implement business logic
	return c.String(http.StatusOK, "Skip action for email "+id)
}
