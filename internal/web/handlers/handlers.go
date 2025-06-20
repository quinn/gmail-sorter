package handlers

import (
	"fmt"

	"github.com/quinn/gmail-sorter/pkg/core"
	"google.golang.org/api/gmail/v1"
)

type Handler struct {
	spec     *core.Spec
	messages []*gmail.Message
}

// NewHandler creates a Handler with a reference to core.Spec
func NewHandler(spec *core.Spec) (*Handler, error) {
	api := spec.GmailService()

	// Fetch the list of messages for navigation
	listRes, err := api.Users.Messages.List("me").MaxResults(50).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	return &Handler{spec: spec, messages: listRes.Messages}, nil
}
