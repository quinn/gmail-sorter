package renderer

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
	"go.quinn.io/ccf/htmx"
)

// EchoRenderer implements the models.Renderer interface and is responsible for
// rendering pages and handling special cases like redirects.
type EchoRenderer struct {
	// Add any necessary fields here
}

// NewEchoRenderer creates a new EchoRenderer.
func NewEchoRenderer() *EchoRenderer {
	return &EchoRenderer{}
}

// Render renders a page based on the action ID and additional data.
func (r *EchoRenderer) Render(c echo.Context, current models.ActionLink, actions []models.ActionLink, data interface{}) error {
	// Render based on actionID
	switch current.Action().ID {
	case "confirm":
		return pages.Confirm(actions).Render(c.Request().Context(), c.Response().Writer)
	case "menu":
		return pages.Menu().Render(c.Request().Context(), c.Response().Writer)
	case "success":
		link, ok := data.(models.ActionLink)
		if !ok {
			return fmt.Errorf("expected ActionLink, got %T", data)
		}
		return pages.Success(actions, link).Render(c.Request().Context(), c.Response().Writer)
	}

	switch data := data.(type) {
	case models.Open:
		return htmx.Redirect(c, data.URL)
	case models.EmailResponse:
		// Used for both /email and /email-group
		if current.Action().ID == "email-group" {
			return pages.GroupEmail(data, actions).Render(c.Request().Context(), c.Response().Writer)
		}
		return pages.Email(data, actions).Render(c.Request().Context(), c.Response().Writer)
	case models.FiltersPageData:
		return pages.Filters(data.AccountID, data.Filters, actions).Render(c.Request().Context(), c.Response().Writer)
	case models.GroupByPageData:
		return pages.GroupBy(data.GroupType, data.Value, data.Emails, actions).Render(c.Request().Context(), c.Response().Writer)
	case models.GroupByDeleteSuccessPageData:
		return pages.GroupByDeleteSuccess(actions, data.GroupType, data.Value, data.Count).Render(c.Request().Context(), c.Response().Writer)
	default:
		return fmt.Errorf("no renderer found for action ID: %s", current.Action().ID)
	}
}
