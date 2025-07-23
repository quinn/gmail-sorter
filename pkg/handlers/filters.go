package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/pkg/db"
)

func init() {
	models.Register(FiltersAction)
}

var FiltersAction models.Action = models.Action{
	ID:      "filters",
	Method:  "GET",
	Path:    "/filters",
	Handler: filters,
	Label:   filtersLabel,
}

func filtersLabel(link models.ActionLink) string {
	return "Filters"
}

func filters(c models.Context) error {
	gmail := middleware.GetGmail(c)

	idStr := c.QueryParam("id")
	if idStr == "" {
		for id := range gmail.API {
			return c.Redirect(FiltersAction.Link(models.WithParams(fmt.Sprintf("%d", id))))
		}
		return errors.New("no accounts")
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	accountID := uint(idInt)

	filters, err := db.DB.AllFilters(accountID)
	if err != nil {
		return err
	}

	actions := []models.ActionLink{
		FiltersRefreshAction.Link(),
	}

	return c.Render(actions, models.FiltersPageData{
		AccountID: accountID,
		Filters:   filters,
	})
}
