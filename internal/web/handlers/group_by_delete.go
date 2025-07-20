package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
	"google.golang.org/api/gmail/v1"
)

func init() {
	models.Register(GroupByDeleteAction)
}

var GroupByDeleteAction models.Action = models.Action{
	ID:               "group-by-delete",
	Method:           "POST",
	Path:             "/account/:id/group-by/:type/delete",
	UnwrappedHandler: groupByDelete,
	Label:            groupByDeleteLabel,
}

func groupByDeleteLabel(link models.ActionLink) string {
	return "Delete By \"" + link.Params[1] + "\": " + link.Fields["val"]
}

func getID(c echo.Context) (uint, error) {
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

func getAPI(c echo.Context) (*gmailapi.GmailAPI, error) {
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

func groupByDelete(c echo.Context) error {
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

	u := "/emails/group-by/" + c.Param("type") + "/delete/success"
	u += "?val=" + url.QueryEscape(c.FormValue("val"))
	u += "&count=" + strconv.Itoa(count)

	return c.Redirect(http.StatusSeeOther, u)
}
