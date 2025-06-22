package models

import (
	"strings"

	"google.golang.org/api/gmail/v1"
)

type EmailResponse struct {
	ID         string
	ThreadID   string
	From       string
	FromDomain string
	To         string
	Subject    string
	Date       string
	Snippet    string
}

func FromGmailMessage(msg *gmail.Message) EmailResponse {
	var from, to, subject, date string

	if msg.Payload == nil {
		panic("FromGmailMessage expects a FullMessage (no payload)")
	}

	for _, h := range msg.Payload.Headers {
		switch h.Name {
		case "From":
			from = h.Value
		case "To":
			to = h.Value
		case "Subject":
			subject = h.Value
		case "Date":
			date = h.Value
		}
	}

	return EmailResponse{
		ID:         msg.Id,
		ThreadID:   msg.ThreadId,
		From:       from,
		FromDomain: from[strings.Index(from, "@")+1:],
		To:         to,
		Subject:    subject,
		Date:       date,
		Snippet:    msg.Snippet,
	}
}
