package middleware

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const dbKey = "gormdb"

// DB injects a *gorm.DB into the echo.Context
func DB(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(dbKey, db)
			return next(c)
		}
	}
}

// GetDB retrieves the *gorm.DB from echo.Context
func GetDB(c echo.Context) *gorm.DB {
	if dbi, ok := c.Get(dbKey).(*gorm.DB); ok {
		return dbi
	}
	return nil
}
