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
	ID:               "group-by-delete",
	Method:           "POST",
	Path:             "/emails/group-by/delete",
	UnwrappedHandler: groupByDelete,
	Label:            groupByDeleteLabel,
}

func groupByDeleteLabel(link models.ActionLink) string {
	return "Delete"
}

func groupByDelete(c echo.Context) error {
	return c.String(http.StatusOK, "Delete action for group by")
}
