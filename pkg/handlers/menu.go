package handlers

import (
	"github.com/quinn/gmail-sorter/internal/web/models"
)

func init() {
	models.Register(MenuAction)
}

var MenuAction models.Action = models.Action{
	ID:      "menu",
	Method:  "GET",
	Path:    "/menu",
	Handler: menu,
	Label:   menuLabel,
}

func menuLabel(link models.ActionLink) string {
	return "Menu"
}

func menu(c models.Context) error {
	return c.Render(nil, nil)
}
