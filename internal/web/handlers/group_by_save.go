package handlers

// func init() {
// 	models.Register(GroupBySaveAction)
// }

// var GroupBySaveAction models.Action = models.Action{
// 	ID:               "group-by-save",
// 	Method:           "GET",
// 	Path:             "/account/:id/group-by/:type/save",
// 	UnwrappedHandler: groupBySave,
// 	Label:            groupBySaveLabel,
// }

// func groupBySaveLabel(link models.ActionLink) string {
// 	return "Group By " + link.Params[0]
// }

// func groupBySave(c echo.Context) error {
// 	groupType := c.Param("type") // domain, from, to
// 	val := c.QueryParam("val")
// 	query, err := groupQuery(c)
// 	if err != nil {
// 		return err
// 	}

// 	// Fetch emails matching the query using Gmail API
// 	api, err := getAPI(c)
// 	if err != nil {
// 		return err
// 	}
// 	slog.Info("Fetching emails matching query: ", "query", query)
// 	res, err := api.Service.Users.Messages.List("me").Q(query).MaxResults(500).Do()
// 	if err != nil {
// 		return echo.NewHTTPError(500, "Failed to fetch emails: "+err.Error())
// 	}

// 	var groupedEmails []models.EmailResponse
// 	for _, m := range res.Messages {
// 		fullMsg, err := api.FullMessage(m.Id)
// 		if err != nil {
// 			return err
// 		}
// 		res, err := models.FromGmailMessage(fullMsg)
// 		if err != nil {
// 			return err
// 		}
// 		groupedEmails = append(groupedEmails, res)
// 	}

// 	actions := []models.ActionLink{
// 		GroupByDeleteAction.Link(
// 			models.WithParams(groupType),
// 			models.WithFields(map[string]string{"val": val}),
// 			models.WithConfirm(),
// 		),
// 	}
// 	return pages.GroupBy(groupType, val, groupedEmails, actions).Render(c.Request().Context(), c.Response().Writer)
// }
