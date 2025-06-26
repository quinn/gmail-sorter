package middleware

import (
	"github.com/labstack/echo/v4"
	"google.golang.org/api/gmail/v1"
)

func Gmail(gmail *gmail.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("gmail", gmail)

			return next(c)
		}
	}
}

func GetGmail(c echo.Context) *gmail.Service {
	return c.Get("gmail").(*gmail.Service)
}
