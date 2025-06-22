package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Index renders the index page
func (h *Handler) Index(c echo.Context) error {

	m := h.messages[0]

	return c.Redirect(http.StatusFound, "/emails/"+m.Id)
	// return pages.Index().Render(c.Request().Context(), c.Response().Writer)
}
