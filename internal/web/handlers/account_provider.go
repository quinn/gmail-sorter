package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

// AccountProviderSelect handles GET /accounts/new
func AccountProviderSelect(c echo.Context) error {
	return pages.ProviderSelect().Render(c.Request().Context(), c.Response().Writer)
}
