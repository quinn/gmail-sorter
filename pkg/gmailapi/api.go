package gmailapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/pkg/db"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailAPI struct {
	Service  *gmail.Service
	DB       *db.DB
	Account  *db.OAuthAccount
	Messages *[]*gmail.Message
}

func (g *GmailAPI) Archive(id string) error {
	moveToArchive := &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"INBOX"},
	}

	if err := g.Modify(id, moveToArchive); err != nil {
		return err
	}

	g.Skip(id)

	return nil
}

// Start is bullshit
func Start(dbConn *db.DB, acct *db.OAuthAccount) (*GmailAPI, error) {
	config, err := models.LoadOauthConfig("google")
	if err != nil {
		return nil, fmt.Errorf("failed to load oauth config: %w", err)
	}

	var token oauth2.Token
	err = json.Unmarshal([]byte(acct.TokenJSON), &token)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	client := config.Client(context.Background(), &token)
	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Gmail client: %w", err)
	}

	a := &GmailAPI{Service: service, DB: dbConn, Account: acct}
	a.RefreshMessages()

	return a, nil
}

func (g *GmailAPI) RefreshMessages() error {
	listRes, err := g.Service.Users.Messages.List("me").Q("in:inbox").MaxResults(50).Do()
	if err != nil {
		return fmt.Errorf("failed to list messages: %w", err)
	}

	g.Messages = &listRes.Messages
	return nil
}

// CreateFilterForGroupDelete creates a Gmail filter for the specified group type and value.
func (g *GmailAPI) CreateFilterForGroupDelete(groupType, val string) error {
	if val == "" {
		return fmt.Errorf("missing val")
	}

	var criteria *gmail.FilterCriteria
	switch groupType {
	case "domain":
		criteria = &gmail.FilterCriteria{From: "*@" + val}
	case "from":
		criteria = &gmail.FilterCriteria{From: val}
	case "to":
		criteria = &gmail.FilterCriteria{To: val}
	default:
		return fmt.Errorf("invalid group type: %s", groupType)
	}

	filterSpec := &gmail.Filter{
		Criteria: criteria,
		Action: &gmail.FilterAction{
			RemoveLabelIds: []string{"INBOX"},
			AddLabelIds:    []string{"TRASH"},
		},
	}

	if err := g.findOrCreateFilter(filterSpec); err != nil {
		return err
	}

	return nil
}

func (g *GmailAPI) FullMessage(id string) (*gmail.Message, error) {
	return g.Service.Users.Messages.Get("me", id).Format("full").Do()
}

func (g *GmailAPI) Filters() ([]*gmail.Filter, error) {
	filters, err := g.DB.GetAll("filters")

	if err != nil {
		return nil, fmt.Errorf("failed to get all filters: %w", err)
	}

	var result []*gmail.Filter

	for _, bytes := range filters {
		var filter gmail.Filter
		if err := json.Unmarshal(bytes, &filter); err != nil {
			return nil, fmt.Errorf("failed to unmarshal filter: %w", err)
		}

		result = append(result, &filter)
	}

	return result, nil
}

func (g *GmailAPI) RefreshFilters() error {
	filters, err := g.Service.Users.Settings.Filters.List("me").Do()
	if err != nil {
		return err
	}

	for _, filter := range filters.Filter {
		d, err := json.Marshal(filter)
		if err != nil {
			return err
		}

		g.DB.Upsert("filters", filter.Id, d)
	}

	labels, err := g.Service.Users.Labels.List("me").Do()
	if err != nil {
		return err
	}

	for _, label := range labels.Labels {
		d, err := json.Marshal(label)
		if err != nil {
			return err
		}

		g.DB.Upsert("labels", label.Id, d)
	}

	return nil
}

func (g *GmailAPI) Label(id string) (*gmail.Label, error) {
	bytes, err := g.DB.Get("labels", id)

	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return nil, fmt.Errorf("label %s not found", id)
	}

	var label gmail.Label
	if err = json.Unmarshal(bytes, &label); err != nil {
		return nil, err
	}

	return &label, nil
}

func (g *GmailAPI) ApplyBatch(query string, batch *gmail.BatchModifyMessagesRequest) (int, error) {
	var res *gmail.ListMessagesResponse

	res, err := g.Service.Users.Messages.List("me").
		MaxResults(500).
		Q(query).
		Do()

	if err != nil {
		return 0, fmt.Errorf("failed to list messages: %w", err)
	}

	var ids []string
	for _, message := range res.Messages {
		ids = append(ids, message.Id)
	}

	batch.Ids = ids

	if err = g.Service.Users.Messages.BatchModify("me", batch).Do(); err != nil {
		return 0, fmt.Errorf("failed to batch modify: %w", err)
	}

	for _, id := range ids {
		g.Skip(id)
	}

	return len(ids), nil
}

func (g *GmailAPI) Delete(id string) error {
	moveToTrash := &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"INBOX"},
		AddLabelIds:    []string{"TRASH"},
	}

	if err := g.Modify(id, moveToTrash); err != nil {
		return err
	}

	g.Skip(id)

	return nil
}

func (g *GmailAPI) Skip(id string) {
	newMessages := (*g.Messages)[:0]
	for _, m := range *g.Messages {
		if m.Id != id {
			newMessages = append(newMessages, m)
		} else {
			slog.Debug("found skip id")
		}
	}
	*g.Messages = newMessages
}

func (g *GmailAPI) Modify(id string, mod *gmail.ModifyMessageRequest) error {
	_, err := g.Service.Users.Messages.Modify("me", id, mod).Do()
	return err
}

func (g *GmailAPI) Query(query string) (*gmail.ListMessagesResponse, error) {
	return g.Service.Users.Messages.List("me").Q(query).MaxResults(500).Do()
}

func (g *GmailAPI) findOrCreateFilter(filterSpec *gmail.Filter) error {
	filters, err := g.DB.GetAll("filters")
	if err != nil {
		return fmt.Errorf("failed to get all filters: %w", err)
	}

	for _, bytes := range filters {
		var filter gmail.Filter
		err = json.Unmarshal(bytes, &filter)

		if err != nil {
			return fmt.Errorf("failed to unmarshal: %w", err)
		}

		if filtersEqual(&filter, filterSpec) {
			slog.Debug("matched filter!")
			return nil
		}
	}

	filter, err := g.Service.Users.Settings.Filters.Create("me", filterSpec).Do()
	if err != nil {
		return fmt.Errorf("could not create filter: %w\n%v", err, filterSpec)
	}

	d, err := json.Marshal(filter)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	if err := g.DB.Upsert("filters", filter.Id, d); err != nil {
		return fmt.Errorf("failed to upsert: %w", err)
	}

	return nil
}
