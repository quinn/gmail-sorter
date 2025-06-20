package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
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
	if title == "" {
		return "Gmail Sorter"
	}
	return title + " | Gmail Sorter"
}
