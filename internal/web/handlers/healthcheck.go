package handlers

import "github.com/labstack/echo/v4"

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(200, map[string]string{"status": "ok"})
}
