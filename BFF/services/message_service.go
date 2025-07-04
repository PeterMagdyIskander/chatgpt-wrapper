package services

import (
	"bff/models"
	"strings"
	"sync"
)

// MessageService handles message-related business logic
type MessageService struct {
	messages []models.UserMessage
	mu       sync.Mutex
}

// NewMessageService creates a new message service
func NewMessageService() *MessageService {
	return &MessageService{
		messages: make([]models.UserMessage, 0),
	}
}

// AddMessage adds a new message to the store
func (s *MessageService) AddMessage(msg models.UserMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, msg)
}

// GetAllMessages returns all messages (thread-safe copy)
func (s *MessageService) GetAllMessages() []models.UserMessage {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]models.UserMessage{}, s.messages...)
}

// GetMessagesByUserId returns messages for a specific user
func (s *MessageService) GetMessagesByUserId(userId string) []models.UserMessage {
	s.mu.Lock()
	defer s.mu.Unlock()

	var userMessages []models.UserMessage
	for _, msg := range s.messages {
		if msg.UserId == userId {
			userMessages = append(userMessages, msg)
		}
	}
	return userMessages
}

// GetMessageCount returns the total number of messages
func (s *MessageService) GetMessageCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.messages)
}

// GetMessageById returns a message by its ID
func (s *MessageService) GetMessageById(id string) (*models.UserMessage, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, msg := range s.messages {
		if msg.ID == id {
			return &msg, true
		}
	}
	return nil, false
}

// GetMessagesContainingKeywords returns messages that contain any of the provided keywords
func (s *MessageService) GetMessagesContainingKeywords(keywords []string) []models.UserMessage {
	s.mu.Lock()
	defer s.mu.Unlock()

	var flaggedMessages []models.UserMessage
	keywordMap := make(map[string]bool)

	// Create a map for faster lookup
	for _, keyword := range keywords {
		keywordMap[keyword] = true
	}

	for _, msg := range s.messages {
		// This is a simple check - in practice you'd want to use the lemmatizer
		// to properly check if the message contains any keywords
		words := strings.Fields(strings.ToLower(msg.Message))
		for _, word := range words {
			cleanWord := strings.Trim(word, ".,!?;:\"'")
			if keywordMap[cleanWord] {
				flaggedMessages = append(flaggedMessages, msg)
				break
			}
		}
	}

	return flaggedMessages
}
