package models

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// EchoContext is a wrapper around echo.Context that implements our custom Context interface.
type EchoContext struct {
	echo.Context
	renderer Renderer
}

// NewEchoContext creates a new EchoContext with the specified renderer.
func NewEchoContext(c echo.Context, renderer Renderer) *EchoContext {
	return &EchoContext{
		Context:  c,
		renderer: renderer,
	}
}

// Redirect redirects to the URL generated from the ActionLink.
func (c *EchoContext) Redirect(link ActionLink) error {
	// Get the action from the link
	action := link.Action()

	// Build URL with appropriate parameters
	url := action.Path

	// Add query parameters if needed - simplified implementation
	// In a real implementation, you might need to handle path parameters as well

	// For now, redirect to the path defined in the action
	return c.Context.Redirect(http.StatusSeeOther, url)
}

// Render renders a template with the provided data and action links.
// - For email page: a gmailapi.Message view
// - For other pages: appropriate data structures needed by those pages
func (c *EchoContext) Render(actions []ActionLink, data any) error {
	link := c.Get("link").(ActionLink)
	return c.renderer.Render(c.Context, link, actions, data)
}
