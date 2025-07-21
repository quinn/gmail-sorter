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
	Service *gmail.Service
	Account *db.OAuthAccount
	// MessageList *MessageList
}

func (g *GmailAPI) Archive(id string) error {
	moveToArchive := &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"INBOX"},
	}

	if err := g.Modify(id, moveToArchive); err != nil {
		return err
	}

	return nil
}

// Start is bullshit
func Start(acct *db.OAuthAccount) (*GmailAPI, error) {
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

	a := &GmailAPI{Service: service, Account: acct}
	return a, nil
}

func (g *GmailAPI) RefreshMessages() ([]*gmail.Message, error) {
	listRes, err := g.Service.Users.Messages.List("me").Q("in:inbox").MaxResults(50).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	var messages []*gmail.Message
	for _, msg := range listRes.Messages {
		fullMsg, err := g.FullMessage(msg.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get email: %w", err)
		}
		messages = append(messages, fullMsg)
	}

	return messages, nil
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

func (g *GmailAPI) RefreshFilters() error {
	filters, err := g.Service.Users.Settings.Filters.List("me").Do()
	if err != nil {
		return err
	}

	if err := db.DB.UpsertFilters(g.Account.ID, filters.Filter); err != nil {
		return err
	}

	labels, err := g.Service.Users.Labels.List("me").Do()
	if err != nil {
		return err
	}

	if err := db.DB.UpsertLabels(g.Account.ID, labels.Labels); err != nil {
		return err
	}

	return nil
}

func (g *GmailAPI) ApplyBatch(query string, batch *gmail.BatchModifyMessagesRequest) ([]string, error) {
	var res *gmail.ListMessagesResponse

	res, err := g.Service.Users.Messages.List("me").
		MaxResults(500).
		Q(query).
		Do()

	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	var ids []string
	for _, message := range res.Messages {
		ids = append(ids, message.Id)
	}

	if len(ids) == 0 {
		return nil, nil
	}

	batch.Ids = ids

	if err = g.Service.Users.Messages.BatchModify("me", batch).Do(); err != nil {
		return nil, fmt.Errorf("failed to batch modify: %w", err)
	}

	return ids, nil
}

func (g *GmailAPI) Delete(id string) error {
	moveToTrash := &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"INBOX"},
		AddLabelIds:    []string{"TRASH"},
	}

	if err := g.Modify(id, moveToTrash); err != nil {
		return err
	}

	return nil
}

func (g *GmailAPI) Modify(id string, mod *gmail.ModifyMessageRequest) error {
	_, err := g.Service.Users.Messages.Modify("me", id, mod).Do()
	return err
}

func (g *GmailAPI) Query(query string) (*gmail.ListMessagesResponse, error) {
	return g.Service.Users.Messages.List("me").Q(query).MaxResults(500).Do()
}

func (g *GmailAPI) findOrCreateFilter(filterSpec *gmail.Filter) error {
	filters, err := db.DB.AllFilters(g.Account.ID)
	if err != nil {
		return err
	}

	for _, filter := range filters {
		if filtersEqual(filter, filterSpec) {
			slog.Debug("matched filter!")
			return nil
		}
	}

	filter, err := g.Service.Users.Settings.Filters.Create("me", filterSpec).Do()
	if err != nil {
		return fmt.Errorf("could not create filter: %w\n%v", err, filterSpec)
	}

	if err := db.DB.UpsertFilter(g.Account.ID, filter); err != nil {
		return fmt.Errorf("failed to upsert filter: %w", err)
	}

	return nil
}
