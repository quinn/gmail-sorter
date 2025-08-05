package tui

import (
	"fmt"
	"strings"

	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/pkg/db"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
	"gorm.io/gorm"
)

// Context implements models.Context for terminal usage.
// For now, URL/query/param helpers return empty strings.
type Context struct {
	page *Page
}

// NewContext constructs a terminal context tied to a Bubble Tea model.
func NewContext(p *Page) *Context {
	return &Context{page: p}
}

func ParamNames(pattern string) ([]string, error) {
	if len(pattern) == 0 || pattern[0] != '/' {
		return nil, fmt.Errorf("invalid pattern %q: must start with '/'", pattern)
	}

	var names []string
	for i, seg := range strings.Split(pattern, "/") {
		if i == 0 || seg == "" {
			continue
		}

		switch {
		case seg[0] == ':':
			if len(seg) == 1 {
				return nil, fmt.Errorf("empty parameter name in %q", pattern)
			}
			names = append(names, seg[1:]) // drop the leading ':'
		}
	}
	return names, nil
}

func (c *Context) Param(name string) string {
	names, err := ParamNames(c.page.Current.Action().Path)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	for idx, n := range names {
		if n == name {
			return c.page.Current.Params[idx]
		}
	}
	return "PARAM_NOT_FOUND"
}

func (c *Context) QueryParam(name string) string { return c.page.Current.Fields[name] }
func (c *Context) FormValue(name string) string  { return c.page.Current.Fields[name] }

// Redirect executes the target action immediately.
func (c *Context) Redirect(link models.ActionLink) error {
	c.page.Current = link
	return c.page.Current.Action().Handler(c)
}

func (c *Context) Render(actions []models.ActionLink, data any) error {
	// Update the model's page; rendering will be handled in Model.View.
	c.page.Actions = actions
	c.page.Data = data
	return nil
}

var gmail *gmailapi.MessageList

func (c *Context) Get(key string) any {
	switch key {
	case "gmail":
		if gmail == nil {
			accts, err := db.DB.GetOAuthAccountsByProvider("google")
			if err == gorm.ErrRecordNotFound {
				panic("no oauth accounts found")
			}
			if err != nil {
				panic(err)
			}
			gmail, err = gmailapi.New(accts)
			if err != nil {
				panic(err)
			}
		}
		return gmail
	}
	return nil
}
