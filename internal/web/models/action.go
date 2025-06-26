package models

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
)

type Action struct {
	Method           string `json:"method"`
	Path             string `json:"path"`
	Label            string `json:"label"`
	Shortcut         string `json:"shortcut"`
	UnwrappedHandler func(c echo.Context) error
}

type ActionContextKey struct{}

func (a Action) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("action", a)
		ctx := c.Request().Context()
		withValue := context.WithValue(ctx, ActionContextKey{}, a)
		c.SetRequest(c.Request().WithContext(withValue))
		return a.UnwrappedHandler(c)
	}
}

func GetAction(ctx context.Context) Action {
	return ctx.Value(ActionContextKey{}).(Action)
}

func (a Action) ID() string {
	return a.Method + " " + a.Path
}

type LinkOpt func(*ActionLink)

func WithConfirm() LinkOpt {
	return func(l *ActionLink) {
		l.Confirm = true
	}
}

func WithPath(path string) LinkOpt {
	return func(l *ActionLink) {
		l.Path = path
	}
}

func WithFields(fields map[string]string) LinkOpt {
	return func(l *ActionLink) {
		l.Fields = fields
	}
}

func (a Action) Link(opts ...LinkOpt) ActionLink {
	link := ActionLink{
		ActionID: a.ID(),
	}
	for _, opt := range opts {
		opt(&link)
	}
	return link
}

type ActionLink struct {
	ActionID string            `json:"action_id"`
	Path     string            `json:"path"`
	Fields   map[string]string `json:"fields,omitempty"`
	Confirm  bool              `json:"confirm,omitempty"`
}

func (l ActionLink) Action() Action {
	parts := strings.Split(l.ActionID, " ")
	if len(parts) != 2 {
		panic("invalid action ID")
	}

	method := parts[0]
	path := parts[1]

	for _, action := range Actions {
		if action.Method == method && action.Path == path {
			return action
		}
	}

	panic("action not found")
}
