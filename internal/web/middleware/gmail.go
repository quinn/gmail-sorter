package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/pkg/db"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
	"gorm.io/gorm"
)

var gmail *gmailapi.MessageList

func Gmail() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if gmail == nil {
				if strings.HasPrefix(c.Path(), "/accounts") ||
					strings.HasPrefix(c.Path(), "/oauth") ||
					strings.HasPrefix(c.Path(), "/healthz") ||
					strings.HasPrefix(c.Path(), "/public") {
					return next(c)
				}

				accts, err := db.DB.GetOAuthAccountsByProvider("google")
				if err == gorm.ErrRecordNotFound {
					return c.Redirect(http.StatusSeeOther, "/accounts")
				}
				if err != nil {
					return err
				}
				gmail, err = gmailapi.New(accts)
				if err != nil {
					return err
				}
			}

			c.Set("gmail", gmail)

			return next(c)
		}
	}
}

type ContextGetter interface {
	Get(key string) any
}

func GetGmail(c ContextGetter) *gmailapi.MessageList {
	val, ok := c.Get("gmail").(*gmailapi.MessageList)
	if !ok {
		return nil
	}

	return val
}
