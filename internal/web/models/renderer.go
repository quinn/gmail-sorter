package models

import (
	"github.com/labstack/echo/v4"
)

// Renderer is responsible for rendering pages based on action IDs and handling
// special data types like Open for redirects.
type Renderer interface {
	// RenderPage renders a page based on the action ID and additional data.
	// The ctx is the request context, w is the response writer,
	// actionID identifies which page to render, and data contains any
	// additional information needed for rendering.
	Render(c echo.Context, current ActionLink, actions []ActionLink, data any) error
}
