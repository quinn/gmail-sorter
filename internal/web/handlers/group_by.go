package handlers

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
)

func init() {
	models.Register(GroupByEmailAction)
}

var GroupByEmailAction models.Action = models.Action{
	ID:               "group-by",
	Method:           "GET",
	Path:             "/account/:id/group-by/:type",
	UnwrappedHandler: groupByEmail,
	Label:            groupByLabel,
}

func groupByLabel(link models.ActionLink) string {
	return "Group By " + link.Params[1]
}

func groupQuery(c echo.Context) (string, error) {
	groupType := c.Param("type") // domain, from, to
	val := c.QueryParam("val")

	if val == "" {
		val = c.FormValue("val")
	}

	if val == "" {
		return "", fmt.Errorf("missing val")
	}

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
		return "", fmt.Errorf("invalid group type: %s", groupType)
	}

	query = "in:inbox " + query
	return query, nil
}

// GroupByEmail handles GET /emails/:id/group/by/:type
func groupByEmail(c echo.Context) error {
	groupType := c.Param("type") // domain, from, to
	val := c.QueryParam("val")
	query, err := groupQuery(c)
	if err != nil {
		return err
	}
	api, err := getAPI(c)
	if err != nil {
		return err
	}

	slog.Info("Fetching emails matching query: ", "query", query)
	res, err := api.Service.Users.Messages.List("me").Q(query).MaxResults(500).Do()
	if err != nil {
		return echo.NewHTTPError(500, "Failed to fetch emails: "+err.Error())
	}

	var groupedEmails []models.EmailResponse
	for _, m := range res.Messages {
		fullMsg, err := api.FullMessage(m.Id)
		if err != nil {
			return err
		}
		res, err := models.FromGmailMessage(fullMsg)
		if err != nil {
			return err
		}
		groupedEmails = append(groupedEmails, res)
	}

	actions := []models.ActionLink{
		GroupByDeleteAction.Link(
			models.WithParams(c.Param("id"), groupType),
			models.WithFields(map[string]string{"val": val}),
			models.WithConfirm(),
		),
	}
	return pages.GroupBy(groupType, val, groupedEmails, actions).Render(c.Request().Context(), c.Response().Writer)
}
