package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/quinn/gmail-sorter/internal/web/models"
)

func init() {
	models.Register(SuccessAction)
}

var SuccessAction models.Action = models.Action{
	ID:      "success",
	Method:  "GET",
	Path:    "/success",
	Handler: success,
	Label:   successLabel,
}

func successLabel(link models.ActionLink) string {
	return "Success"
}

func success(c models.Context) error {
	link := c.QueryParam("link")
	var actionObj models.ActionLink
	if err := json.Unmarshal([]byte(link), &actionObj); err != nil {
		return fmt.Errorf("failed to unmarshal action: %w", err)
	}

	actions := []models.ActionLink{
		IndexAction.Link(),
	}

	return c.Render(actions, actionObj)
}
