package models

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

// EchoContext is a wrapper around echo.Context that implements our custom Context interface.
type EchoContext struct {
	echo.Context
}

// NewEchoContext creates a new EchoContext.
func NewEchoContext(c echo.Context) *EchoContext {
	return &EchoContext{Context: c}
}

// FormValue returns the value from the POST form.
func (c *EchoContext) FormValue(name string) string {
	return c.Context.FormValue(name)
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
	// Get the request context and response writer once
	ctx := c.Request().Context()
	writer := c.Response().Writer

	switch link.Action().ID {
	case "confirm":
		return pages.Confirm(actions).Render(ctx, writer)
	}

	// Use type switching to determine what to render based on the data type
	switch data.(type) {
	// no cases yet, but we will need this.

	default:
		return fmt.Errorf("cannot render data of type %T: no appropriate renderer found", data)
	}
}
