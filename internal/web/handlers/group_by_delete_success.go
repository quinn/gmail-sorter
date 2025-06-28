package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

// GroupByDeleteSuccessAction is the action for the post-delete success screen
var GroupByDeleteSuccessAction = models.Action{
	ID:               "group-by-delete-success",
	Method:           "GET",
	Path:             "/emails/group-by/:type/delete/success",
	UnwrappedHandler: groupByDeleteSuccess,
	Label:            groupByDeleteSuccessLabel,
}

func init() {
	models.Register(GroupByDeleteSuccessAction)
}

func groupByDeleteSuccessLabel(link models.ActionLink) string {
	return "Post Delete Success"
}

func groupByDeleteSuccess(c echo.Context) error {
	groupType := c.Param("type")
	val := c.QueryParam("val")
	countStr := c.QueryParam("count")
	count, _ := strconv.Atoi(countStr)

	actions := []models.ActionLink{
		GroupBySaveFilterAction.Link(
			models.WithParams(groupType),
			models.WithFields(map[string]string{"val": val}),
		),
		IndexAction.Link(),
	}

	return pages.GroupByDeleteSuccess(actions, groupType, val, count).Render(c.Request().Context(), c.Response().Writer)
}
