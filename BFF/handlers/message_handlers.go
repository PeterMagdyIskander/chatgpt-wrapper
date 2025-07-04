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
	var newMessage models.UserMessage

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

	// Generate ID if not provided
	if newMessage.ID == "" {
		newMessage.ID = generateMessageID()
	}

	// Add message to database first
	h.messageService.AddMessage(newMessage)

	// Check if message contains any forbidden keywords
	foundKeywords := h.keywordService.CheckTextForKeywords(newMessage.Message)

	if len(foundKeywords) > 0 {
		// Message contains forbidden keywords - return 400
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Message contains forbidden keywords",
			"messageId":     newMessage.ID,
			"foundKeywords": foundKeywords,
			"message":       "Your message has been saved but contains prohibited content",
		})
		return
	}

	// Message is clean - return 200
	c.JSON(http.StatusOK, gin.H{
		"messageId": newMessage.ID,
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

// GetMessagesByUser handles GET /messages/user/:userId
func (h *MessageHandlers) GetMessagesByUser(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserId parameter is required"})
		return
	}

	messages := h.messageService.GetMessagesByUserId(userId)
	c.JSON(http.StatusOK, messages)
}

// GetMessageStats handles GET /messages/stats
func (h *MessageHandlers) GetMessageStats(c *gin.Context) {
	count := h.messageService.GetMessageCount()
	c.JSON(http.StatusOK, gin.H{
		"total_messages": count,
	})
}

// GetMessageById handles GET /messages/:id
func (h *MessageHandlers) GetMessageById(c *gin.Context) {
	messageId := c.Param("id")
	if messageId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message ID parameter is required"})
		return
	}

	message, found := h.messageService.GetMessageById(messageId)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.JSON(http.StatusOK, message)
}
