package gmailapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"reflect"

	"github.com/pkg/errors"
	"github.com/quinn/gmail-sorter/pkg/db"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v2"
)

var tokFile = "token.json"

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokFile)

	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(token)

	if err != nil {
		log.Fatalf("Could not encode JSON: %v", err)
	}
}

type GmailAPI struct {
	Service  *gmail.Service
	db       *db.DB
	Messages *[]*gmail.Message
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

	_, err := g.findOrCreateFilter(filterSpec)
	return err
}

// Start is bullshit
func Start(db *db.DB) (*GmailAPI, error) {
	b, err := os.ReadFile("credentials.json")

	if err != nil {
		return nil, errors.Errorf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b,
		gmail.GmailModifyScope,
		gmail.GmailSettingsBasicScope,
	)

	if err != nil {
		return nil, errors.Errorf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)
	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		return nil, errors.Errorf("Unable to retrieve Gmail client: %v", err)
	}

	// Fetch the list of messages for navigation
	listRes, err := service.Users.Messages.List("me").MaxResults(50).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	return &GmailAPI{Service: service, db: db, Messages: &listRes.Messages}, nil
}

func (g *GmailAPI) FullMessage(id string) (*gmail.Message, error) {
	return g.Service.Users.Messages.Get("me", id).Format("full").Do()
}

func (g *GmailAPI) ApplyFilter(filterSpec *gmail.Filter, query string, batch gmail.BatchModifyMessagesRequest) error {
	_, err := g.findOrCreateFilter(filterSpec)
	if err != nil {
		return err
	}

	pageToken := ""

	for {
		var res *gmail.ListMessagesResponse

		res, err = g.Service.Users.Messages.List("me").
			MaxResults(50).
			PageToken(pageToken).
			Q(query).
			Do()

		if err != nil {
			return fmt.Errorf("failed to list messages: %v", err)
		}

		var ids []string
		for _, message := range res.Messages {
			ids = append(ids, message.Id)
		}

		if len(ids) == 0 {
			break
		}

		pageToken = res.NextPageToken

		batch.Ids = ids

		err = g.Service.Users.Messages.BatchModify("me", &batch).Do()

		if err != nil {
			return fmt.Errorf("failed to batch modify: %v", err)
		}

		if pageToken == "" {
			break
		} else {
			slog.Info("continuing to next page", "query", query)
		}
	}

	return nil
}

func (g *GmailAPI) Filters() ([]*gmail.Filter, error) {
	filters, err := g.db.GetAll("filters")

	if err != nil {
		return nil, fmt.Errorf("failed to get all filters: %v", err)
	}

	var result []*gmail.Filter

	for _, bytes := range filters {
		var filter gmail.Filter
		err = yaml.Unmarshal(bytes, &filter)

		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %v", err)
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
		var d []byte
		d, err = yaml.Marshal(filter)

		if err != nil {
			return err
		}

		g.db.Upsert("filters", filter.Id, d)
	}

	labels, err := g.Service.Users.Labels.List("me").Do()
	if err != nil {
		return err
	}

	for _, label := range labels.Labels {
		var d []byte
		d, err = yaml.Marshal(label)

		if err != nil {
			return err
		}

		g.db.Upsert("labels", label.Id, d)
	}

	return nil
}

func (g *GmailAPI) Label(id string) (*gmail.Label, error) {
	bytes, err := g.db.Get("labels", id)

	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return nil, fmt.Errorf("label %s not found", id)
	}

	var label gmail.Label
	if err = yaml.Unmarshal(bytes, &label); err != nil {
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
		return 0, fmt.Errorf("failed to list messages: %v", err)
	}

	var ids []string
	for _, message := range res.Messages {
		ids = append(ids, message.Id)
	}

	batch.Ids = ids

	if err = g.Service.Users.Messages.BatchModify("me", batch).Do(); err != nil {
		return 0, fmt.Errorf("failed to batch modify: %v", err)
	}

	return len(ids), nil
}

func (g *GmailAPI) Query(query string) (*gmail.ListMessagesResponse, error) {
	return g.Service.Users.Messages.List("me").Q(query).MaxResults(500).Do()
}

func (g *GmailAPI) findOrCreateFilter(filterSpec *gmail.Filter) (_ *gmail.Filter, err error) {
	filters, err := g.db.GetAll("filters")

	if err != nil {
		return nil, fmt.Errorf("failed to get all filters: %v", err)
	}

	for _, bytes := range filters {
		var filter gmail.Filter
		err = yaml.Unmarshal(bytes, &filter)

		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %v", err)
		}

		match := false

		if filterSpec.Action.AddLabelIds == nil {
			filterSpec.Action.AddLabelIds = []string{}
		}

		if filterSpec.Action.RemoveLabelIds == nil {
			filterSpec.Action.RemoveLabelIds = []string{}
		}

		if reflect.DeepEqual(filter.Action.AddLabelIds, filterSpec.Action.AddLabelIds) &&
			reflect.DeepEqual(filter.Action.RemoveLabelIds, filterSpec.Action.RemoveLabelIds) {
			match = true
		}

		if match {
			slog.Debug("matched filter!")
			return
		}
	}

	filter, err := g.Service.Users.Settings.Filters.Create("me", filterSpec).Do()

	if err != nil {
		err = errors.Errorf("could not create filter: %v\n%v", err, filterSpec)
		return
	}

	d, err := yaml.Marshal(filter)
	if err != nil {
		return
	}

	err = g.db.Upsert("filters", filter.Id, d)

	if err != nil {
		return
	}

	return
}
