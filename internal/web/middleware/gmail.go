package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/pkg/db"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
	"gorm.io/gorm"
)

var gmail *gmailapi.GmailAPI

func Gmail(db *db.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if gmail == nil {
				if strings.HasPrefix(c.Path(), "/accounts") ||
					strings.HasPrefix(c.Path(), "/oauth") ||
					strings.HasPrefix(c.Path(), "/healthz") ||
					strings.HasPrefix(c.Path(), "/public") {
					return next(c)
				}

				acct, err := db.GetOAuthAccountByProvider("google")
				if err == gorm.ErrRecordNotFound {
					return c.Redirect(http.StatusSeeOther, "/accounts")
				}
				if err != nil {
					return err
				}

				gmail, err = gmailapi.Start(db, acct)
				if err != nil {
					return err
				}
			}

			c.Set("gmail", gmail)

			return next(c)
		}
	}
}

func GetGmail(c echo.Context) *gmailapi.GmailAPI {
	val, ok := c.Get("gmail").(*gmailapi.GmailAPI)
	if !ok {
		return nil
	}

	return val
}
