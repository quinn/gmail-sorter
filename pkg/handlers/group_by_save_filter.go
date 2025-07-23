package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/quinn/gmail-sorter/internal/web/models"
)

// GroupBySaveFilterAction is the action for the Gmail filter creation screen
var GroupBySaveFilterAction = models.Action{
	ID:      "group-by-save-filter",
	Method:  "POST",
	Path:    "/account/:id/group-by/:type/save-filter",
	Handler: groupBySaveFilter,
	Label:   groupBySaveFilterLabel,
}

func init() {
	models.Register(GroupBySaveFilterAction)
}

func groupBySaveFilterLabel(link models.ActionLink) string {
	return "Save as Gmail Filter"
}

func groupBySaveFilter(c models.Context) error {
	groupType := c.Param("type")
	val := c.FormValue("val")
	api, err := getAPI(c)
	if err != nil {
		return err
	}

	if err := api.CreateFilterForGroupDelete(groupType, val); err != nil {
		return fmt.Errorf("failed to create gmail filter: %w", err)
	}

	linkJSON, err := json.Marshal(c.Get("link"))
	if err != nil {
		return fmt.Errorf("failed to marshal link: %w", err)
	}

	link := SuccessAction.Link(
		models.WithFields(map[string]string{"link": string(linkJSON)}),
	)
	return c.Redirect(link)
}
