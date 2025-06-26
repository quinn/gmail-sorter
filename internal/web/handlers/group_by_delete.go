package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

func init() {
	models.Register(GroupByDeleteAction)
}

var GroupByDeleteAction models.Action = models.Action{
	Method:           "POST",
	Path:             "/emails/group-by/delete",
	Label:            "Delete",
	Shortcut:         "d",
	UnwrappedHandler: groupByDelete,
}

func groupByDelete(c echo.Context) error {
	return c.String(http.StatusOK, "Delete action for group by")
}
