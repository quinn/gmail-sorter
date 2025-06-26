package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

func init() {
	models.Register(MenuAction)
}

var MenuAction models.Action = models.Action{
	ID:               "menu",
	Method:           "GET",
	Path:             "/menu",
	UnwrappedHandler: menu,
	Label:            menuLabel,
}

func menuLabel(link models.ActionLink) string {
	return "Menu"
}

func menu(c echo.Context) error {
	return pages.Menu().Render(c.Request().Context(), c.Response().Writer)
}
