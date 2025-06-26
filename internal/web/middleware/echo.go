package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
)

type EchoContextKey struct{}

func Echo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		newCtx := context.WithValue(ctx, EchoContextKey{}, c)
		c.SetRequest(c.Request().WithContext(newCtx))

		return next(c)
	}
}

func GetEcho(ctx context.Context) echo.Context {
	return ctx.Value(EchoContextKey{}).(echo.Context)
}
