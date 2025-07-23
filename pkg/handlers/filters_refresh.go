package handlers

import (
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

func init() {
	FiltersRefreshAction = models.Action{
		ID:      "filters-refresh",
		Method:  "POST",
		Path:    "/filters/refresh",
		Handler: filtersRefresh,
		Label:   filtersRefreshLabel,
	}
	models.Register(FiltersRefreshAction)
}

var FiltersRefreshAction models.Action


func filtersRefreshLabel(link models.ActionLink) string {
	return "Refresh Filters"
}

func filtersRefresh(c models.Context) error {
	api := middleware.GetGmail(c)

	if err := api.RefreshFilters(); err != nil {
		return err
	}

	return c.Redirect(FiltersAction.Link())
}
