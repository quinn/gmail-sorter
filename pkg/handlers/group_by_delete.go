package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
	"google.golang.org/api/gmail/v1"
)

func init() {
	models.Register(GroupByDeleteAction)
}

var GroupByDeleteAction models.Action = models.Action{
	ID:      "group-by-delete",
	Method:  "POST",
	Path:    "/account/:id/group-by/:type/delete",
	Handler: groupByDelete,
	Label:   groupByDeleteLabel,
}

func groupByDeleteLabel(link models.ActionLink) string {
	return "Delete By \"" + link.Params[1] + "\": " + link.Fields["val"]
}

func getID(c models.Context) (uint, error) {
	idStr := c.Param("id")
	var accountID uint
	if idStr == "" {
		return 0, errors.New("missing account id")
	}

	if idInt, err := strconv.Atoi(idStr); err != nil {
		return 0, err
	} else {
		accountID = uint(idInt)
	}

	return accountID, nil
}

func getAPI(c models.Context) (*gmailapi.GmailAPI, error) {
	gm := middleware.GetGmail(c)
	accountID, err := getID(c)
	if err != nil {
		return nil, err
	}

	api, ok := gm.API[accountID]
	if !ok {
		return nil, fmt.Errorf("account API %d not found", accountID)
	}

	return api, nil
}

func groupByDelete(c models.Context) error {
	gm := middleware.GetGmail(c)
	accountID, err := getID(c)
	if err != nil {
		return err
	}

	query, err := groupQuery(c)
	if err != nil {
		return err
	}

	batch := gmail.BatchModifyMessagesRequest{
		RemoveLabelIds: []string{"INBOX"},
		AddLabelIds:    []string{"TRASH"},
	}

	count, err := gm.ApplyBatch(accountID, query, &batch)
	if err != nil {
		return err
	}

	link := GroupByDeleteSuccessAction.Link(
		models.WithParams(c.Param("id"), c.Param("type")),
		models.WithFields(map[string]string{"val": c.FormValue("val"), "count": strconv.Itoa(count)}),
	)

	return c.Redirect(link)
}
