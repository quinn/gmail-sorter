package util

import (
	"fmt"

	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
)

func GetEmail(c middleware.ContextGetter, id string) (*gmailapi.Message, error) {
	gm := middleware.GetGmail(c)

	email, err := gm.GetFullMessage(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get email: %w", err)
	}

	return email, nil
}
