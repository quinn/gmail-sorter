package models

type Context interface {
	Param(name string) string
	QueryParam(name string) string
	FormValue(name string) string
	Redirect(link ActionLink) error
	Render(actions []ActionLink, data any) error
	Get(key string) interface{}
}

type Handler func(c Context) error
