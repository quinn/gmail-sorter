package handlers

import (
	"strings"
	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

// GroupByEmail handles GET /emails/:id/group/by/:type
func (h *Handler) GroupByEmail(c echo.Context) error {
	id := c.Param("id")
	groupType := c.Param("type") // domain, from, to

	msg, _, err := h.getEmail(id)
	if err != nil {
		return err
	}

	// Build Gmail search query based on groupType
	var query string
	switch groupType {
	case "domain":
		// Extract domain from the sender's email
		from := ""
		for _, h := range msg.Payload.Headers {
			if h.Name == "From" {
				from = h.Value
				break
			}
		}
		domain := ""
		if at := len(from); at > 0 {
			parts := strings.Split(from, "@")
			if len(parts) == 2 {
				domain = parts[1]
			}
		}
		if domain == "" {
			return echo.NewHTTPError(400, "Could not extract domain from sender")
		}
		query = "from:*@" + domain
	case "from":
		from := ""
		for _, h := range msg.Payload.Headers {
			if h.Name == "From" {
				from = h.Value
				break
			}
		}
		if from == "" {
			return echo.NewHTTPError(400, "Could not extract sender from email")
		}
		query = "from:'" + from + "'"
	case "to":
		to := ""
		for _, h := range msg.Payload.Headers {
			if h.Name == "To" {
				to = h.Value
				break
			}
		}
		if to == "" {
			return echo.NewHTTPError(400, "Could not extract recipient from email")
		}
		query = "to:'" + to + "'"
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
