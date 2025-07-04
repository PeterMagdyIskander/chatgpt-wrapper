package handlers

import (
	"bff/models"
	"bff/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// MessageHandlers contains all message-related HTTP handlers
type MessageHandlers struct {
	messageService *services.MessageService
	keywordService *services.KeywordService
}

// NewMessageHandlers creates a new message handlers instance
func NewMessageHandlers(messageService *services.MessageService, keywordService *services.KeywordService) *MessageHandlers {
	return &MessageHandlers{
		messageService: messageService,
		keywordService: keywordService,
	}
}

// PostMessage handles POST /messages
func (h *MessageHandlers) PostMessage(c *gin.Context) {
	var newMessage models.UserMessageDTO

	if err := c.BindJSON(&newMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Basic validation
	if newMessage.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message cannot be empty"})
		return
	}

	if newMessage.UserId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserId cannot be empty"})
		return
	}

	// Check if message contains any forbidden keywords
	foundKeywords := h.keywordService.CheckTextForKeywords(newMessage.Message)

	message := models.MessageUserTable{
		MessageId:      generateMessageID(),
		UserId:         newMessage.UserId,
		Status:         len(foundKeywords) > 0,
		MessageContent: newMessage.Message,
	}

	// Add message to database first
	h.messageService.AddMessage(message)

	if len(foundKeywords) > 0 {
		// Message contains forbidden keywords - return 400
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Message contains forbidden keywords",
			"messageId":     message.MessageId,
			"foundKeywords": foundKeywords,
			"message":       "Your message has been saved but contains prohibited content",
		})
		return
	}

	// Message is clean - return 200
	c.JSON(http.StatusOK, gin.H{
		"messageId": message.MessageId,
		"message":   "Message posted successfully",
		"status":    "approved",
	})
}

// generateMessageID creates a simple message ID
// In production, you might want to use UUID or a more sophisticated ID generator
func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

// GetMessages handles GET /messages
func (h *MessageHandlers) GetMessages(c *gin.Context) {
	messages := h.messageService.GetAllMessages()
	c.JSON(http.StatusOK, messages)
}
