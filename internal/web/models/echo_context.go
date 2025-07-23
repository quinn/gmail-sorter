package models

import (
	"github.com/labstack/echo/v4"
)

// EchoContext is a wrapper around echo.Context that implements our custom Context interface.
type EchoContext struct {
	echo.Context
}

// NewEchoContext creates a new EchoContext.
func NewEchoContext(c echo.Context) *EchoContext {
	return &EchoContext{Context: c}
}
