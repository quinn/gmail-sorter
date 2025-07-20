package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
)

func Messages() echo.MiddlewareFunc {
	list := gmailapi.MessageList{}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("messages", &list)
			return next(c)
		}
	}
}

func GetMessages(c echo.Context) *gmailapi.MessageList {
	return c.Get("messages").(*gmailapi.MessageList)
}
