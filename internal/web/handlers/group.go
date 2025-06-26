package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/util"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

func init() {
	models.Register(EmailGroupAction)
}

var EmailGroupAction models.Action = models.Action{
	Method:           "GET",
	Path:             "/emails/:id/group",
	Label:            "Group",
	Shortcut:         "g",
	UnwrappedHandler: groupEmail,
}

// GroupEmail handles GET /emails/:id/group
func groupEmail(c echo.Context) error {
	id := c.Param("id")

	email, err := util.GetEmail(c, id)
	if err != nil {
		return err
	}

	actions := []models.ActionLink{
		GroupByEmailAction.Link(
			models.WithPath("/emails/group-by/domain"),
			models.WithFields(map[string]string{"val": email.FromDomain}),
		),
		GroupByEmailAction.Link(
			models.WithPath("/emails/group-by/from"),
			models.WithFields(map[string]string{"val": email.From}),
		),
		GroupByEmailAction.Link(
			models.WithPath("/emails/group-by/to"),
			models.WithFields(map[string]string{"val": email.To}),
		),
	}

	return pages.GroupEmail(email, actions).Render(c.Request().Context(), c.Response().Writer)
}
