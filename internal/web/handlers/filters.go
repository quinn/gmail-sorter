package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
	"github.com/quinn/gmail-sorter/pkg/db"
)

func init() {
	models.Register(FiltersAction)
}

var FiltersAction models.Action = models.Action{
	ID:               "filters",
	Method:           "GET",
	Path:             "/filters",
	UnwrappedHandler: filters,
	Label:            filtersLabel,
}

func filtersLabel(link models.ActionLink) string {
	return "Filters"
}

func filters(c echo.Context) error {
	gmail := middleware.GetGmail(c)

	idStr := c.QueryParam("id")
	if idStr == "" {
		for id, _ := range gmail.API {
			return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/filters?id=%d", id))
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

	// filters := []*gmail.Filter{}

	actions := []models.ActionLink{
		FiltersRefreshAction.Link(),
	}

	return pages.Filters(accountID, filters, actions).Render(c.Request().Context(), c.Response().Writer)
}
