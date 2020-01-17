package core

import (
	"fmt"
	"io/ioutil"
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
	Domains []string `yaml:"domains"`
	api     *gmail.Service
	db      *db.DB
}

const timeFormat = "Mon Jan 2 15:04:05 -0700 MST 2006"
const testlabelname = "stackoverflow.email"

func NewSpec(api *gmail.Service, db *db.DB) (spec Spec, err error) {
	log.Info("Starting new spec")
	bytes, err := ioutil.ReadFile("./spec.yaml")
	spec.api = api
	spec.db = db

	if err != nil {
		return
	}

	err = yaml.Unmarshal(bytes, &spec)

	if err != nil {
		return
	}

	return
}

func (s *Spec) Apply() (err error) {
	oneMinuteAgo := time.Now().Add(time.Duration(-3) * time.Hour)
	refreshTimestamp, err := s.getTimestamp("refresh")

	if err != nil {
		return
	}

	if refreshTimestamp == nil || refreshTimestamp.Before(oneMinuteAgo) {
		err = s.refreshLabels()

		if err != nil {
			return
		}

		err = s.refreshFilters()

		if err != nil {
			return
		}

		_, err = s.setTimestamp("refresh")

		if err != nil {
			return
		}
	}

	return s.createDomainFilter(testlabelname)
}

func (s *Spec) refreshLabels() (err error) {
	log.Info("refreshing labels")
	r, err := s.api.Users.Labels.List("me").Do()

	if err != nil {
		return
	}

	for _, label := range r.Labels {
		var d []byte
		d, err = yaml.Marshal(label)

		if err != nil {
			return
		}

		err = s.db.Upsert("labels", label.Id, d)

		if err != nil {
			return
		}
	}

	return
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

func (s *Spec) findOrCreateLabel(labelName string) (label *gmail.Label, err error) {
	labels, err := s.db.GetAll("labels")

	if err != nil {
		return
	}

	for _, bytes := range labels {
		var l gmail.Label
		err = yaml.Unmarshal(bytes, &l)
		label = &l

		if err != nil {
			return
		}

		if label.Name == labelName {
			return
		}
	}

	label = &gmail.Label{
		Name:                  labelName,
		LabelListVisibility:   "labelShow",
		MessageListVisibility: "show",
	}

	label, err = s.api.Users.Labels.Create("me", label).Do()

	if err != nil {
		return
	}

	d, err := yaml.Marshal(label)

	if err != nil {
		return
	}

	err = s.db.Upsert("labels", label.Id, d)

	if err != nil {
		return
	}

	return
}

func (s *Spec) findOrCreateFilter(label *gmail.Label, filterSpec *gmail.Filter) (filter *gmail.Filter, err error) {
	filters, err := s.db.GetAll("labels")

	if err != nil {
		return
	}

	for _, bytes := range filters {
		var f gmail.Filter
		err = yaml.Unmarshal(bytes, &f)
		filter = &f

		if err != nil {
			return
		}

		if filter.Action != nil && filter.Action.AddLabelIds[0] == label.Id {
			log.Info("Matched label!")
			return
		}
	}

	filter, err = s.api.Users.Settings.Filters.Create("me", filterSpec).Do()

	if err != nil {
		err = errors.Errorf("Could not create filter: %s", err)
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

func (s *Spec) createDomainFilter(domain string) (err error) {
	labelName := fmt.Sprintf("Domains/%s", domain)
	label, err := s.findOrCreateLabel(labelName)

	if err != nil {
		return
	}

	filterSpec := gmail.Filter{
		Action: &gmail.FilterAction{
			AddLabelIds: []string{label.Id},
		},
		Criteria: &gmail.FilterCriteria{
			From: domain,
		},
	}

	_, err = s.findOrCreateFilter(label, &filterSpec)

	if err != nil {
		return
	}

	query := fmt.Sprintf("from:%s", domain)
	res, err := s.api.Users.Messages.List("me").
		MaxResults(5).
		Q(query).
		Do()

	if err != nil {
		return
	}

	spew.Dump(res)

	s.api.Users.Messages.BatchModify("me", &gmail.BatchModifyMessagesRequest{
		AddLabelIds: []string{label.Id},
	})

	return
}

func (s *Spec) getTimestamp(name string) (_ *time.Time, err error) {
	bytes, err := s.db.Get("timestamps", name)

	if err != nil {
		return
	}

	if bytes == nil {
		return
	}

	timestamp, err := time.Parse(timeFormat, string(bytes))

	if err != nil {
		return
	}

	return &timestamp, nil
}

func (s *Spec) setTimestamp(name string) (now time.Time, err error) {
	now = time.Now()
	str := now.Format(timeFormat)
	err = s.db.Upsert("timestamps", name, []byte(str))

	if err != nil {
		return
	}

	return
}
