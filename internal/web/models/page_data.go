package models

import "google.golang.org/api/gmail/v1"

// FiltersPageData is the data struct passed to the filters page renderer
// Used for type-safe rendering from handlers and renderer.
type FiltersPageData struct {
	AccountID uint
	Filters   []*gmail.Filter
}

// GroupByPageData is used for rendering the group-by page.
type GroupByPageData struct {
	GroupType string
	Value     string
	Emails    []EmailResponse
}

// GroupByDeleteSuccessPageData is used for rendering the group-by delete success page.
type GroupByDeleteSuccessPageData struct {
	GroupType string
	Value     string
	Count     int
}
