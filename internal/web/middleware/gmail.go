package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
)

func Gmail(gmail *gmailapi.GmailAPI) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("gmail", gmail)

			return next(c)
		}
	}
}

func GetGmail(c echo.Context) *gmailapi.GmailAPI {
	return c.Get("gmail").(*gmailapi.GmailAPI)
}
