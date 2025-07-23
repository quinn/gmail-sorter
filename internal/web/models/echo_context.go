package models

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func LinkFormAction(c echo.Context, link ActionLink) string {
	var params []any
	for _, param := range link.Params {
		params = append(params, param)
	}

	return c.Echo().Reverse(link.ActionID, params...)
}

func LinkURL(c echo.Context, link ActionLink) (string, error) {
	if link.Action().Method != "GET" {
		return "", fmt.Errorf("method not GET")
	}

	a := LinkFormAction(c, link)
	u, err := url.Parse(a)
	if err != nil {
		return "", fmt.Errorf("failed to parse echo reversed url: %w", err)
	}

	qs := u.Query()
	for name, value := range link.Fields {
		qs.Set(name, value)
	}
	u.RawQuery = qs.Encode()
	return u.String(), nil
}

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
	url, err := LinkURL(c.Context, link)
	if err != nil {
		return err
	}

	return c.Context.Redirect(http.StatusSeeOther, url)
}

// Render renders a template with the provided data and action links.
// - For email page: a gmailapi.Message view
// - For other pages: appropriate data structures needed by those pages
func (c *EchoContext) Render(actions []ActionLink, data any) error {
	link := c.Get("link").(ActionLink)
	return c.renderer.Render(c.Context, link, actions, data)
}
