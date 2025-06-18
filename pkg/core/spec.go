package core

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/quinn/gmail-sorter/pkg/db"
	"google.golang.org/api/gmail/v1"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

// Spec represents the spec.yaml
type Spec struct {
	Domains     []string `yaml:"domains"`
	Delete      []string `yaml:"delete"`
	Newsletters []string `yaml:"newsletters"`
	api         *gmail.Service
	db          *db.DB
}

const timeFormat = "Mon Jan 2 15:04:05 -0700 MST 2006"

// NewSpec loads the spec.yaml
func NewSpec(api *gmail.Service, db *db.DB) (*Spec, error) {
	log.SetLevel(log.DebugLevel)
	log.Info("starting new spec")

	bytes, err := os.ReadFile("./spec.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	spec := &Spec{
		api: api,
		db:  db,
	}

	err = yaml.Unmarshal(bytes, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}

	return spec, nil
}

// Apply Creates labels and filters for the spec.yaml
func (s *Spec) Apply() error {
	timeAgo := time.Now().Add(time.Duration(-3) * time.Hour)
	refreshTimestamp, err := s.getTimestamp("refresh")
	if err != nil {
		return fmt.Errorf("failed to get timestamp: %v", err)
	}

	if refreshTimestamp == nil || refreshTimestamp.Before(timeAgo) {
		log.Debug("refresh time has expired. refreshing")

		err = s.refreshLabels()
		if err != nil {
			return err
		}

		err = s.refreshFilters()
		if err != nil {
			return err
		}

		_, err = s.setTimestamp("refresh")
		if err != nil {
			return err
		}
	}

	for _, domain := range s.Domains {
		if err := s.createDomainFilter(domain); err != nil {
			return err
		}
	}

	for _, newsletterDomain := range s.Newsletters {
		if err := s.createNewsletterFilter(newsletterDomain); err != nil {
			return err
		}
	}

	return nil
}

func (s *Spec) refreshLabels() error {
	log.Info("refreshing labels")
	r, err := s.api.Users.Labels.List("me").Do()

	if err != nil {
		return fmt.Errorf("failed to list labels: %v", err)
	}

	for _, label := range r.Labels {
		var d []byte
		d, err = yaml.Marshal(label)

		if err != nil {
			return err
		}

		err = s.db.Upsert("labels", label.Id, d)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Spec) refreshFilters() (err error) {
	r, err := s.api.Users.Settings.Filters.List("me").Do()

	if err != nil {
		return
	}

	for _, filter := range r.Filter {
		var d []byte
		d, err = yaml.Marshal(filter)

		if err != nil {
			return
		}

		spew.Dump(filter.Id)
		err = s.db.Upsert("filters", filter.Id, d)

		if err != nil {
			return
		}
	}

	return
}

func (s *Spec) findOrCreateLabel(labelName string) (*gmail.Label, error) {
	labels, err := s.db.GetAll("labels")
	if err != nil {
		return nil, fmt.Errorf("failed db getAll labels: %v", err)
	}

	for _, bytes := range labels {
		var label gmail.Label
		err = yaml.Unmarshal(bytes, &label)

		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal label: %v", err)
		}

		if strings.EqualFold(label.Name, labelName) {
			return &label, nil
		}
	}

	label := &gmail.Label{
		Name:                  labelName,
		LabelListVisibility:   "labelShow",
		MessageListVisibility: "show",
	}

	createdLabel, err := s.api.Users.Labels.Create("me", label).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create label \"%s\" in gmail: %v", label.Name, err)
	} else {
		label = createdLabel
	}

	d, err := yaml.Marshal(label)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}

	err = s.db.Upsert("labels", label.Id, d)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert: %v", err)
	}

	return label, nil
}

func (s *Spec) findOrCreateFilter(filterSpec *gmail.Filter) (_ *gmail.Filter, err error) {
	filters, err := s.db.GetAll("filters")

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
			log.Debug("matched filter!")
			return
		}
	}

	filter, err := s.api.Users.Settings.Filters.Create("me", filterSpec).Do()

	if err != nil {
		err = errors.Errorf("could not create filter: %v\n%v", err, filterSpec)
		return
	}

	d, err := yaml.Marshal(filter)
	if err != nil {
		return
	}

	err = s.db.Upsert("filters", filter.Id, d)

	if err != nil {
		return
	}

	return
}

func (s *Spec) createNewsletterFilter(domain string) error {
	labelName := "Newsletters"
	label, err := s.findOrCreateLabel(labelName)
	if err != nil {
		return err
	}

	filterSpec := gmail.Filter{
		Action: &gmail.FilterAction{
			AddLabelIds:    []string{label.Id},
			RemoveLabelIds: []string{"INBOX"},
		},
		Criteria: &gmail.FilterCriteria{
			From: domain,
		},
	}

	query := fmt.Sprintf("from:%s in:inbox", domain)

	batch := gmail.BatchModifyMessagesRequest{
		AddLabelIds:    []string{label.Id},
		RemoveLabelIds: []string{"INBOX"},
	}

	return s.applyFilter(&filterSpec, query, batch)
}

func (s *Spec) createDomainFilter(domain string) error {
	labelName := fmt.Sprintf("Domains/%s", domain)
	label, err := s.findOrCreateLabel(labelName)
	if err != nil {
		return err
	}

	filterSpec := gmail.Filter{
		Action: &gmail.FilterAction{
			AddLabelIds: []string{label.Id},
		},
		Criteria: &gmail.FilterCriteria{
			From: domain,
		},
	}

	query := fmt.Sprintf("from:%s in:inbox", domain)

	batch := gmail.BatchModifyMessagesRequest{
		AddLabelIds: []string{label.Id},
	}

	return s.applyFilter(&filterSpec, query, batch)
}

func (s *Spec) applyFilter(filterSpec *gmail.Filter, query string, batch gmail.BatchModifyMessagesRequest) error {
	_, err := s.findOrCreateFilter(filterSpec)
	if err != nil {
		return err
	}

	pageToken := ""

	for {
		var res *gmail.ListMessagesResponse

		res, err = s.api.Users.Messages.List("me").
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

		err = s.api.Users.Messages.BatchModify("me", &batch).Do()

		if err != nil {
			return fmt.Errorf("failed to batch modify: %v", err)
		}

		if pageToken == "" {
			break
		} else {
			log.Debugf("continuing to next page: \"%s\"", query)
		}
	}

	return nil
}

func (s *Spec) getTimestamp(name string) (*time.Time, error) {
	bytes, err := s.db.Get("timestamps", name)

	if err != nil {
		return nil, fmt.Errorf("failed db get: %v", err)
	}

	if bytes == nil {
		return nil, nil
	}

	timestamp, err := time.Parse(timeFormat, string(bytes))

	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}

	return &timestamp, nil
}

func (s *Spec) setTimestamp(name string) (time.Time, error) {
	now := time.Now()
	str := now.Format(timeFormat)
	err := s.db.Upsert("timestamps", name, []byte(str))
	if err != nil {
		return time.Time{}, fmt.Errorf("failed db upsert: %v", err)
	}

	return now, nil
}
