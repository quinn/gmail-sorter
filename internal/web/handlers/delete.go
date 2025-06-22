package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// DeleteEmail handles POST /emails/:id/delete
func (h *Handler) DeleteEmail(c echo.Context) error {
	id := c.Param("id")
	// TODO: Implement business logic
	return c.String(http.StatusOK, "Delete action for email "+id)
}
