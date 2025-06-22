package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

// GroupByEmail handles GET /emails/:id/group/by/:type
func (h *Handler) GroupByEmail(c echo.Context) error {
	id := c.Param("id")
	groupType := c.Param("type") // domain, from, to

	// Placeholder: fetch all emails in the thread or all messages (for demo)
	var groupedEmails []models.EmailResponse
	for _, m := range h.messages {
		if m.Id == id || m.ThreadId == id {
			groupedEmails = append(groupedEmails, models.FromGmailMessage(m))
		}
	}

	// Placeholder: grouping logic by type
	// In real code, group groupedEmails by domain/from/to
	// For now, just pass all

	return pages.GroupBy(groupType, groupedEmails).Render(c.Request().Context(), c.Response().Writer)
}
