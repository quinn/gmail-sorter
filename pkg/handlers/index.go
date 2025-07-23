package handlers

import (
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

func init() {
	models.Register(IndexAction)
}

var IndexAction models.Action = models.Action{
	ID:      "index",
	Method:  "GET",
	Path:    "/",
	Label:   indexLabel,
	Handler: index,
}

func indexLabel(link models.ActionLink) string {
	return "Home"
}

// index renders the index page
func index(c models.Context) error {
	gm := middleware.GetGmail(c)

	if len(gm.Messages) == 0 {
		if err := gm.Refresh(); err != nil {
			return err
		}
	}

	m := gm.Messages[0]
	link := EmailAction.Link(models.WithParams(m.Message.Id))
	return c.Redirect(link)
}
