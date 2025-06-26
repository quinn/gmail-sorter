package handlers

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

func init() {
	models.Register(GroupBySaveAction)
}

var GroupBySaveAction models.Action = models.Action{
	ID:               "group-by-save",
	Method:           "GET",
	Path:             "/emails/group-by/:type/save",
	UnwrappedHandler: groupBySave,
	Label:            groupBySaveLabel,
}

func groupBySaveLabel(link models.ActionLink) string {
	return "Group By " + link.Params[0]
}

func groupBySave(c echo.Context) error {
	groupType := c.Param("type") // domain, from, to
	val := c.QueryParam("val")
	query, err := groupQuery(c)
	if err != nil {
		return err
	}

	// Fetch emails matching the query using Gmail API
	api := middleware.GetGmail(c)
	slog.Info("Fetching emails matching query: ", "query", query)
	res, err := api.Service.Users.Messages.List("me").Q(query).MaxResults(500).Do()
	if err != nil {
		return echo.NewHTTPError(500, "Failed to fetch emails: "+err.Error())
	}

	var groupedEmails []models.EmailResponse
	for _, m := range res.Messages {
		fullMsg, err := api.FullMessage(m.Id)
		if err != nil {
			continue // skip bad messages
		}
		groupedEmails = append(groupedEmails, models.FromGmailMessage(fullMsg))
	}

	actions := []models.ActionLink{
		GroupByDeleteAction.Link(
			models.WithParams(groupType),
			models.WithFields(map[string]string{"val": val}),
			models.WithConfirm(),
		),
	}
	return pages.GroupBy(groupType, val, groupedEmails, actions).Render(c.Request().Context(), c.Response().Writer)
}
