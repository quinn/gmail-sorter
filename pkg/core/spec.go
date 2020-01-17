package core

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/quinn/gmail-sorter/pkg/db"
	"google.golang.org/api/gmail/v1"
	"gopkg.in/yaml.v2"
)

// Spec represents the spec.yaml
type Spec struct {
	Domains []string `yaml:"domains"`
	api     *gmail.Service
	db      *db.DB
}

const timeFormat = "Mon Jan 2 15:04:05 -0700 MST 2006"

func NewSpec(api *gmail.Service, db *db.DB) (spec Spec, err error) {
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

	return s.createDomainFilter("berniesanders.com")
}

func (s *Spec) refreshLabels() (err error) {
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

		spew.Dump(label.Id)
		err = s.db.Upsert("labels", label.Id, d)

		if err != nil {
			return
		}
	}

	retur
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

func (s *Spec) findOrCreateLabel(labelName string) (label gmail.Label, err error) {
	labels, err := s.db.GetAll("labels")

	if err != nil {
		return
	}

	for _, bytes := range labels {
		err = yaml.Unmarshal(bytes, &label)

		if err != nil {
			return
		}

		if label.Name == labelName {
			return
		}
	}

	label.Name = labelName
	_, err = s.api.Users.Labels.Create("me", &label).Do()

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

	spew.Dump(label)
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
