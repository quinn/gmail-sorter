package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// IndexHandler renders the index page
func (h *Handler) IndexHandler(c echo.Context) error {

	m := h.messages[0]

	return c.Redirect(http.StatusFound, "/emails/"+m.Id)
	// return pages.Index().Render(c.Request().Context(), c.Response().Writer)
}
