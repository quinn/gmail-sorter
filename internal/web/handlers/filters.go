package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
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
	api := middleware.GetGmail(c)

	filters, err := api.Filters()

	if err != nil {
		return err
	}

	actions := []models.ActionLink{
		FiltersRefreshAction.Link(),
	}

	return pages.Filters(api, filters, actions).Render(c.Request().Context(), c.Response().Writer)
}
