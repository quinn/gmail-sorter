package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

type TitleContextKey struct{}

// SetTitle sets the page title in the context
func SetTitle(c echo.Context, title string) {
	ctx := c.Request().Context()
	newCtx := context.WithValue(ctx, TitleContextKey{}, title)
	c.SetRequest(c.Request().WithContext(newCtx))
}

func GetTitle(ctx context.Context) string {
	title, _ := ctx.Value(TitleContextKey{}).(string)
	link := models.GetLink(ctx)
	if link == nil {
		return "Gmail Sorter"
	}
	action := link.Action()
	label := action.Label(*link)
	if label != "" {
		return label + " | Gmail Sorter"
	}
	if title == "" {
		return "Gmail Sorter"
	}
	return title + " | Gmail Sorter"
}
