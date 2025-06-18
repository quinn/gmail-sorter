package handlers

import "github.com/quinn/gmail-sorter/pkg/core"

type Handler struct {
	spec *core.Spec
}

// NewHandler creates a Handler with a reference to core.Spec
func NewHandler(spec *core.Spec) *Handler {
	return &Handler{spec: spec}
}
