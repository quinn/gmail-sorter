package models

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Action struct {
	ID      string `json:"id"`
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler func(c Context) error
	Label   func(link ActionLink) string
}

type LinkContextKey struct{}

func (a Action) WrappedHandler(renderer Renderer) echo.HandlerFunc {
	return func(c echo.Context) error {
		fields := map[string]string{}
		if err := c.Bind(&fields); err != nil {
			return err
		}
		link := a.Link(
			WithParams(c.ParamValues()...),
			WithFields(fields),
		)
		c.Set("link", link)
		ctx := c.Request().Context()
		withValue := context.WithValue(ctx, LinkContextKey{}, &link)
		c.SetRequest(c.Request().WithContext(withValue))
		return a.Handler(NewEchoContext(c, renderer))
	}
}

func GetLink(ctx context.Context) *ActionLink {
	val := ctx.Value(LinkContextKey{})
	if val == nil {
		return nil
	}
	return val.(*ActionLink)
}

type LinkOpt func(*ActionLink)

func WithConfirm() LinkOpt {
	return func(l *ActionLink) {
		l.Confirm = true
	}
}

func WithParams(params ...string) LinkOpt {
	return func(l *ActionLink) {
		l.Params = params
	}
}

func WithFields(fields map[string]string) LinkOpt {
	return func(l *ActionLink) {
		l.Fields = fields
	}
}

func (a Action) Link(opts ...LinkOpt) ActionLink {
	link := ActionLink{
		ActionID: a.ID,
	}
	for _, opt := range opts {
		opt(&link)
	}
	return link
}

type ActionLink struct {
	ActionID string            `json:"action_id"`
	Params   []string          `json:"params,omitempty"`
	Fields   map[string]string `json:"fields,omitempty"`
	Confirm  bool              `json:"confirm,omitempty"`
}

func (l ActionLink) Action() Action {
	for _, action := range Actions {
		if action.ID == l.ActionID {
			return action
		}
	}

	panic("action not found")
}
