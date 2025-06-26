package middleware

import (
	"github.com/labstack/echo/v4"
	"google.golang.org/api/gmail/v1"
)

func Messages(messages *[]*gmail.Message) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("messages", messages)

			return next(c)
		}
	}
}

func GetMessages(c echo.Context) *[]*gmail.Message {
	return c.Get("messages").(*[]*gmail.Message)
}
