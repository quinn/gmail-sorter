package models

var Actions []Action

func Register(action Action) {
	Actions = append(Actions, action)
}
