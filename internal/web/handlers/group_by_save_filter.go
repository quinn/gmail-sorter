package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

// GroupBySaveFilterAction is the action for the Gmail filter creation screen
var GroupBySaveFilterAction = models.Action{
	ID:               "group-by-save-filter",
	Method:           "POST",
	Path:             "/emails/group-by/:type/save-filter",
	UnwrappedHandler: groupBySaveFilter,
	Label:            groupBySaveFilterLabel,
}

func init() {
	models.Register(GroupBySaveFilterAction)
}

func groupBySaveFilterLabel(link models.ActionLink) string {
	return "Save as Gmail Filter"
}

func groupBySaveFilter(c echo.Context) error {
	groupType := c.Param("type")
	val := c.FormValue("val")
	api := middleware.GetGmail(c)

	// This is where Gmail filter creation logic should go
	err := api.CreateFilterForGroupDelete(groupType, val)
	if err != nil {
		return fmt.Errorf("failed to create gmail filter: %w", err)
	}

	return c.Redirect(http.StatusSeeOther, "/")

}
