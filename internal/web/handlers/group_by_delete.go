package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"google.golang.org/api/gmail/v1"
)

func init() {
	models.Register(GroupByDeleteAction)
}

var GroupByDeleteAction models.Action = models.Action{
	ID:               "group-by-delete",
	Method:           "POST",
	Path:             "/emails/group-by/:type/delete",
	UnwrappedHandler: groupByDelete,
	Label:            groupByDeleteLabel,
}

func groupByDeleteLabel(link models.ActionLink) string {
	return "Delete By \"" + link.Params[0] + "\": " + link.Fields["val"]
}

func groupByDelete(c echo.Context) error {
	query, err := groupQuery(c)
	if err != nil {
		return err
	}

	batch := gmail.BatchModifyMessagesRequest{
		RemoveLabelIds: []string{"INBOX"},
		AddLabelIds:    []string{"TRASH"},
	}

	api := middleware.GetGmail(c)
	count, err := api.ApplyBatch(query, &batch)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/emails/group-by/"+c.Param("type")+"/delete/success?val="+c.FormValue("val")+"&count="+strconv.Itoa(count))
}
