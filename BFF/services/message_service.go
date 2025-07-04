package services

import (
	"bff/models"
	"sync"
)

// MessageService handles message-related business logic
type MessageService struct {
	messages []models.MessageUserTable
	mu       sync.Mutex
}

// NewMessageService creates a new message service
func NewMessageService() *MessageService {
	return &MessageService{
		messages: make([]models.MessageUserTable, 0),
	}
}

// AddMessage adds a new message to the store
func (s *MessageService) AddMessage(msg models.MessageUserTable) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, msg)
}

// GetAllMessages returns all messages (thread-safe copy)
func (s *MessageService) GetAllMessages() []models.MessageUserTable {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]models.MessageUserTable{}, s.messages...)
}

// GetMessageById returns a message by its ID
func (s *MessageService) GetMessageById(messageId string) (*models.MessageUserTable, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, msg := range s.messages {
		if msg.MessageId == messageId {
			return &msg, true
		}
	}
	return nil, false
}
