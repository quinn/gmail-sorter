package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

func init() {
	models.Register(FiltersRefreshAction)
}

var FiltersRefreshAction models.Action = models.Action{
	ID:               "filters-refresh",
	Method:           "POST",
	Path:             "/filters/refresh",
	UnwrappedHandler: filtersRefresh,
	Label:            filtersRefreshLabel,
}

func filtersRefreshLabel(link models.ActionLink) string {
	return "Refresh Filters"
}

func filtersRefresh(c echo.Context) error {
	api := middleware.GetGmail(c)

	if err := api.RefreshFilters(); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/filters")
}
