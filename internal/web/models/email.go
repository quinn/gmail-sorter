package models

type EmailResponse struct {
	ID       string
	ThreadID string
	From     string
	To       string
	Subject  string
	Date     string
	Snippet  string
}
