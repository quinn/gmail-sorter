package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/util"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

func init() {
	models.Register(EmailGroupAction)
}

var EmailGroupAction models.Action = models.Action{
	ID:               "email-group",
	Method:           "GET",
	Path:             "/emails/:id/group",
	UnwrappedHandler: groupEmail,
	Label:            groupLabel,
}

func groupLabel(link models.ActionLink) string {
	return "Group"
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
			models.WithParams(strconv.Itoa(int(email.AccountID)), "domain"),
			models.WithFields(map[string]string{"val": email.View.FromDomain}),
		),
	}
	for _, from := range email.View.From {
		actions = append(actions,
			GroupByEmailAction.Link(
				models.WithParams(strconv.Itoa(int(email.AccountID)), "from"),
				models.WithFields(map[string]string{"val": from}),
			))
	}
	for _, to := range email.View.To {
		actions = append(actions,
			GroupByEmailAction.Link(
				models.WithParams(strconv.Itoa(int(email.AccountID)), "to"),
				models.WithFields(map[string]string{"val": to}),
			))
	}

	return pages.GroupEmail(email.View, actions).Render(c.Request().Context(), c.Response().Writer)
}
