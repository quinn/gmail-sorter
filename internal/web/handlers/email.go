package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/util"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

func init() {
	models.Register(EmailAction)
}

var EmailAction models.Action = models.Action{
	Method:           "GET",
	Path:             "/emails/:id",
	Label:            "Email",
	Shortcut:         "e",
	UnwrappedHandler: email,
}

// Email renders a single email by ID
func email(c echo.Context) error {
	id := c.Param("id")
	email, err := util.GetEmail(c, id)
	if err != nil {
		return err
	}

	actions := []models.ActionLink{
		EmailGroupAction.Link(
			models.WithPath("/emails/" + id + "/group"),
		),
		EmailSkipAction.Link(
			models.WithPath("/emails/" + id + "/skip"),
		),
		EmailDeleteAction.Link(
			models.WithPath("/emails/" + id + "/delete"),
		),
	}

	return pages.Email(email, actions).Render(c.Request().Context(), c.Response().Writer)
}
