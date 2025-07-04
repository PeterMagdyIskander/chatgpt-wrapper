package handlers

import (
	"bff/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SSEHandlers contains all Server-Sent Events related HTTP handlers
type SSEHandlers struct {
	messageService *services.MessageService
	openaiService  *services.OpenAIService
}

// NewSSEHandlers creates a new SSE handlers instance
func NewSSEHandlers(messageService *services.MessageService, openaiService *services.OpenAIService) *SSEHandlers {
	return &SSEHandlers{
		messageService: messageService,
		openaiService:  openaiService,
	}
}

// StreamCompletion handles GET /stream with SSE for OpenAI completions
func (h *SSEHandlers) StreamCompletion(c *gin.Context) {
	// Get parameters from query string
	userId := c.Query("userId")
	messageId := c.Query("messageId")

	// Validate required parameters
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId parameter is required"})
		return
	}

	if messageId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "messageId parameter is required"})
		return
	}

	// Check if the message exists and belongs to the user
	message, exists := h.messageService.GetMessageById(messageId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	if message.UserId != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Message does not belong to the specified user"})
		return
	}

	if !message.Status {
		c.JSON(http.StatusForbidden, gin.H{"error": "Message contains forbidden keywords"})
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// Send initial connection message
	c.SSEvent("connection", "Connected to OpenAI stream")
	c.Writer.Flush()

	// Create channels for communication with OpenAI service
	responseChan := make(chan string, 100)
	errorChan := make(chan error, 1)

	// Start OpenAI streaming in a goroutine
	go h.openaiService.StreamCompletion(message.MessageContent, responseChan, errorChan)

	// Handle the streaming response
	for {
		select {
		case content, ok := <-responseChan:
			if !ok {
				// Channel closed, streaming finished
				c.SSEvent("done", "Stream completed")
				c.Writer.Flush()
				return
			}

			// Send the content chunk to the client
			c.SSEvent("data", content)
			c.Writer.Flush()

		case err := <-errorChan:
			if err != nil {
				// Send error to client
				c.SSEvent("error", fmt.Sprintf("Error: %s", err.Error()))
				c.Writer.Flush()
				return
			}

		case <-c.Request.Context().Done():
			// Client disconnected
			return
		}
	}
}
