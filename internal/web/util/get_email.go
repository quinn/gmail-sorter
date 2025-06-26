package util

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"google.golang.org/api/gmail/v1"
)

func GetEmail(c echo.Context, id string) (email models.EmailResponse, err error) {
	messages := middleware.GetMessages(c)
	var msg *gmail.Message
	var idx int
	for i, m := range *messages {
		if m.Id == id {
			msg = m
			idx = i
			break
		}
	}

	if msg == nil {
		return email, fmt.Errorf("email not found")
	}

	if msg.Payload == nil {
		fullMsg, err := middleware.GetGmail(c).Users.Messages.Get("me", msg.Id).Format("full").Do()
		if err != nil {
			return email, fmt.Errorf("failed to get email: %w", err)
		}
		msg = fullMsg
		(*messages)[idx] = fullMsg
	}

	email = models.FromGmailMessage(msg)
	return email, nil
}
