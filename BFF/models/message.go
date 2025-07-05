package models

// UserMessage represents a user message in the system
type CharLimitDTO struct {
	CharLimit int16 `json:"charLimit"`
}

// UserMessage represents a user message in the system
type UserMessageDTO struct {
	Message string `json:"message"`
	UserId  string `json:"userId"`
}

type MessageUserTable struct {
	MessageId      string
	UserId         string
	Flagged        bool
	MessageContent string
}
