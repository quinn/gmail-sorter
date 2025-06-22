package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

// GroupByEmail handles GET /emails/:id/group/by/:type
func (h *Handler) GroupByEmail(c echo.Context) error {
	groupType := c.Param("type") // domain, from, to
	val := c.QueryParam("val")

	// Build Gmail search query based on groupType
	var query string
	switch groupType {
	case "domain":
		query = "from:*@" + val
	case "from":
		query = "from:'" + val + "'"
	case "to":
		query = "to:'" + val + "'"
	default:
		return echo.NewHTTPError(400, "Invalid group type")
	}

	// Fetch emails matching the query using Gmail API
	api := h.spec.GmailService()
	res, err := api.Users.Messages.List("me").Q(query).MaxResults(50).Do()
	if err != nil {
		return echo.NewHTTPError(500, "Failed to fetch emails: "+err.Error())
	}

	var groupedEmails []models.EmailResponse
	for _, m := range res.Messages {
		fullMsg, err := api.Users.Messages.Get("me", m.Id).Format("full").Do()
		if err != nil {
			continue // skip bad messages
		}
		groupedEmails = append(groupedEmails, models.FromGmailMessage(fullMsg))
	}

	return pages.GroupBy(groupType, groupedEmails).Render(c.Request().Context(), c.Response().Writer)
}
