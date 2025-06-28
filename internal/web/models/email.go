package models

import (
	"fmt"
	"net/mail"
	"strings"

	"google.golang.org/api/gmail/v1"
)

type EmailResponse struct {
	ID         string
	ThreadID   string
	From       []string
	FromDomain string
	To         []string
	Subject    string
	Date       string
	Snippet    string
}

func FromGmailMessage(msg *gmail.Message) (EmailResponse, error) {
	var from, to []string
	var subject, date string

	if msg.Payload == nil {
		panic("FromGmailMessage expects a FullMessage (no payload)")
	}

	for _, h := range msg.Payload.Headers {
		switch h.Name {
		case "From":
			addr, err := mail.ParseAddressList(h.Value)
			if err != nil {
				return EmailResponse{}, fmt.Errorf("failed to parse from: %w", err)
			}
			from = make([]string, len(addr))
			for i, a := range addr {
				from[i] = a.Address
			}
		case "To":
			addr, err := mail.ParseAddressList(h.Value)
			if err != nil {
				return EmailResponse{}, fmt.Errorf("failed to parse to: %w", err)
			}
			to = make([]string, len(addr))
			for i, a := range addr {
				to[i] = a.Address
			}
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
		FromDomain: from[0][strings.Index(from[0], "@")+1:],
		To:         to,
		Subject:    subject,
		Date:       date,
		Snippet:    msg.Snippet,
	}, nil
}
