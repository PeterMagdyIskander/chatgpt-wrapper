package services

import (
	"bff/models"
	"sync"
)

type MessageService struct {
	messages  []models.MessageUserTable
	charLimit int16
	mu        sync.Mutex
}

func NewMessageService() *MessageService {
	return &MessageService{
		messages:  make([]models.MessageUserTable, 0),
		charLimit: 100,
	}
}

func (s *MessageService) AddMessage(msg models.MessageUserTable) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, msg)
}

func (s *MessageService) GetAllMessages() []models.MessageUserTable {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]models.MessageUserTable{}, s.messages...)
}

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
func (s *MessageService) SetCharLimit(newCharLimit int16) int16 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.charLimit = newCharLimit
	return s.charLimit
}
func (s *MessageService) GetCharLimit() int16 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.charLimit
}
