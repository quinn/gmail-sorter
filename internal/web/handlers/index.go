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
	Method:           "GET",
	Path:             "/",
	Label:            "Index",
	Shortcut:         "i",
	UnwrappedHandler: index,
}

// index renders the index page
func index(c echo.Context) error {
	m := (*middleware.GetMessages(c))[0]
	return c.Redirect(http.StatusFound, "/emails/"+m.Id)
}
