package models

// UserMessage represents a user message in the system
type UserMessage struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	UserId  string `json:"userId"`
}
