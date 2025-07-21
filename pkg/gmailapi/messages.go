package gmailapi

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/pkg/db"
	"google.golang.org/api/gmail/v1"
)

type Message struct {
	Message   *gmail.Message
	AccountID uint
	View      models.EmailResponse
}

func New(accounts []db.OAuthAccount) (*MessageList, error) {
	var ml MessageList
	ml.API = make(map[uint]*GmailAPI)
	ml.Messages = make([]*Message, 0)
	for _, acct := range accounts {
		g, err := Start(&acct)
		if err != nil {
			return nil, err
		}
		ml.API[acct.ID] = g
	}

	if err := ml.Refresh(); err != nil {
		return nil, err
	}

	return &ml, nil
}

type MessageList struct {
	Messages []*Message
	API      map[uint]*GmailAPI
}

func (m *MessageList) GetFullMessage(id string) (*Message, error) {
	var msg *Message
	var idx int
	for i, m := range m.Messages {
		if m.Message.Id == id {
			msg = m
			idx = i
			break
		}
	}

	if msg == nil {
		return nil, fmt.Errorf("email not found")
	}

	if msg.Message.Payload == nil {
		fullMsg, err := m.API[msg.AccountID].FullMessage(msg.Message.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get email: %w", err)
		}
		msg.Message = fullMsg
		m.Messages[idx].Message = fullMsg
	}

	if msg.View.ID == "" {
		var err error
		if msg.View, err = models.FromGmailMessage(msg.Message); err != nil {
			return nil, fmt.Errorf("failed to convert email: %w", err)
		}
	}

	return msg, nil
}

func (m *MessageList) ApplyBatch(accountID uint, query string, request *gmail.BatchModifyMessagesRequest) (int, error) {
	api, ok := m.API[accountID]
	if !ok {
		return 0, fmt.Errorf("account API %d not found", accountID)
	}

	ids, err := api.ApplyBatch(query, request)
	if err != nil {
		return 0, err
	}

	for _, id := range ids {
		m.Skip(id)
	}

	return len(ids), nil
}

func (m *MessageList) RefreshFilters() error {
	for _, api := range m.API {
		if err := api.RefreshFilters(); err != nil {
			return err
		}
	}

	return nil
}

func (m *MessageList) Archive(id string) error {
	var found bool
	for _, msg := range m.Messages {
		if msg.Message.Id == id {
			found = true
			if err := m.API[msg.AccountID].Archive(id); err != nil {
				return err
			}
		}
	}

	if !found {
		return fmt.Errorf("message %s not found", id)
	}

	m.Skip(id)
	return nil
}

func (m *MessageList) Delete(id string) error {
	var found bool
	for _, msg := range m.Messages {
		if msg.Message.Id == id {
			found = true
			if err := m.API[msg.AccountID].Delete(id); err != nil {
				return err
			}
		}
	}

	if !found {
		return fmt.Errorf("message %s not found", id)
	}

	m.Skip(id)
	return nil
}

func (m *MessageList) OpenURL(id string) (string, error) {
	for _, msg := range m.Messages {
		if msg.Message.Id == id {
			return fmt.Sprintf("https://mail.google.com/mail/u/%d/#inbox/%s", msg.AccountID, id), nil
		}
	}

	return "", fmt.Errorf("message %s not found", id)
}

func (m *MessageList) Refresh() error {
	for _, api := range m.API {
		messages, err := api.RefreshMessages()
		if err != nil {
			return err
		}
		for _, msg := range messages {
			m.Messages = append(m.Messages, &Message{
				Message:   msg,
				AccountID: api.Account.ID,
			})
		}
	}

	m.Sort()
	return nil
}

func (m *MessageList) Skip(id string) {
	newMessages := m.Messages[:0]
	for _, msg := range m.Messages {
		if msg.Message.Id != id {
			newMessages = append(newMessages, msg)
		} else {
			slog.Debug("found skip id")
		}
	}
	m.Messages = newMessages
}

func (m *MessageList) Sort() {
	// Sort by gmail.Message.InternalDate (descending)
	slices.SortFunc(m.Messages, func(a, b *Message) int {
		ma := a.Message
		mb := b.Message

		// InternalDate is int64 (ms since epoch)
		if ma.InternalDate > mb.InternalDate {
			return -1
		}
		if ma.InternalDate < mb.InternalDate {
			return 1
		}
		return 0
	})
}
