package handlers

import (
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/util"
)

var EmailAction models.Action

func init() {
	EmailAction = models.Action{
		ID:      "email",
		Method:  "GET",
		Path:    "/emails/:id",
		Handler: email,
		Label:   emailLabel,
	}
	models.Register(EmailAction)
}

func emailLabel(link models.ActionLink) string {
	return "Email"
}

// Email renders a single email by ID
func email(c models.Context) error {
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

	return c.Render(actions, email.View)
}
