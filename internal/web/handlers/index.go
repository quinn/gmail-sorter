package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

func init() {
	models.Register(IndexAction)
}

var IndexAction models.Action = models.Action{
	ID:               "index",
	Method:           "GET",
	Path:             "/",
	Label:            indexLabel,
	UnwrappedHandler: index,
}

func indexLabel(link models.ActionLink) string {
	return "Home"
}

// index renders the index page
func index(c echo.Context) error {
	api := middleware.GetGmail(c)

	if len(*api.Messages) == 0 {
		if err := api.RefreshMessages(); err != nil {
			return err
		}
	}

	m := (*api.Messages)[0]
	return c.Redirect(http.StatusFound, "/emails/"+m.Id)
}
