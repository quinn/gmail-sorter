package models

type Action struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Label    string `json:"label"`
	Shortcut string `json:"shortcut"`
}
