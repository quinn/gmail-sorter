package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

func init() {
	models.Register(SuccessAction)
}

var SuccessAction models.Action = models.Action{
	ID:               "success",
	Method:           "GET",
	Path:             "/success",
	UnwrappedHandler: success,
	Label:            successLabel,
}

func successLabel(link models.ActionLink) string {
	return "Success"
}

func success(c echo.Context) error {

	link := c.QueryParam("link")
	var actionObj models.ActionLink
	if err := json.Unmarshal([]byte(link), &actionObj); err != nil {
		return fmt.Errorf("failed to unmarshal action: %w", err)
	}

	actions := []models.ActionLink{
		IndexAction.Link(),
	}

	return pages.Success(actions, actionObj).Render(c.Request().Context(), c.Response().Writer)
}
