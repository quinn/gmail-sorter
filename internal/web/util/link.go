package util

import (
	"fmt"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/quinn/gmail-sorter/internal/web/models"
)

func LinkFormAction(c echo.Context, link models.ActionLink) string {
	var params []any
	for _, param := range link.Params {
		params = append(params, param)
	}

	return c.Echo().Reverse(link.ActionID, params...)
}

func LinkURL(c echo.Context, link models.ActionLink) (string, error) {
	if link.Action().Method != "GET" {
		return "", fmt.Errorf("method not GET")
	}

	a := LinkFormAction(c, link)
	u, err := url.Parse(a)
	if err != nil {
		return "", fmt.Errorf("failed to parse echo reversed url: %w", err)
	}

	qs := u.Query()
	for name, value := range link.Fields {
		qs.Set(name, value)
	}
	u.RawQuery = qs.Encode()
	return u.String(), nil
}
