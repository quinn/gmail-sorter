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
	ID:               "email",
	Method:           "GET",
	Path:             "/emails/:id",
	UnwrappedHandler: email,
	Label:            emailLabel,
}

func emailLabel(link models.ActionLink) string {
	return "Email"
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
			models.WithParams(id),
		),
		EmailCommandAction.Link(
			models.WithParams(id, "skip"),
		),
		EmailCommandAction.Link(
			models.WithParams(id, "archive"),
		),
		EmailCommandAction.Link(
			models.WithParams(id, "delete"),
		),
		EmailCommandAction.Link(
			models.WithParams(id, "open"),
		),
		EmailCommandAction.Link(
			models.WithParams(id, "todo"),
		),
	}

	return pages.Email(email.View, actions).Render(c.Request().Context(), c.Response().Writer)
}
