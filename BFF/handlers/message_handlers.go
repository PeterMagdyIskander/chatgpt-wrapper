package handlers

import (
	"bff/models"
	"bff/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


type MessageHandlers struct {
	messageService *services.MessageService
	keywordService *services.KeywordService
}


func NewMessageHandlers(messageService *services.MessageService, keywordService *services.KeywordService) *MessageHandlers {
	return &MessageHandlers{
		messageService: messageService,
		keywordService: keywordService,
	}
}


func (h *MessageHandlers) PostMessage(c *gin.Context) {
	var newMessage models.UserMessageDTO

	if err := c.BindJSON(&newMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if newMessage.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message cannot be empty"})
		return
	}

	if len(newMessage.Message) > int(h.messageService.GetCharLimit()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Character Size"})
		return
	}

	if newMessage.UserId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserId cannot be empty"})
		return
	}


	foundKeywords := h.keywordService.CheckTextForKeywords(newMessage.Message)

	message := models.MessageUserTable{
		MessageId:      generateMessageID(),
		UserId:         newMessage.UserId,
		Flagged:        len(foundKeywords) > 0,
		MessageContent: newMessage.Message,
	}


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


	c.JSON(http.StatusOK, gin.H{
		"messageId": message.MessageId,
		"message":   "Message posted successfully",
		"status":    "approved",
	})
}



func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}


func (h *MessageHandlers) GetMessages(c *gin.Context) {
	messages := h.messageService.GetAllMessages()
	c.JSON(http.StatusOK, messages)
}


func (h *MessageHandlers) PostCharLimit(c *gin.Context) {
	var newVarLimit models.CharLimitDTO

	// Bind the JSON request body to newVarLimit
	if err := c.ShouldBindJSON(&newVarLimit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("CharLimit received:", newVarLimit.CharLimit)



	charLimit := h.messageService.SetCharLimit(newVarLimit.CharLimit)
	res := models.CharLimitDTO{
		CharLimit: charLimit,
	}
	c.JSON(http.StatusOK, res)
}



func (h *MessageHandlers) GetCharLimit(c *gin.Context) {
	charLimit := h.messageService.GetCharLimit()
	res := models.CharLimitDTO{
		CharLimit: charLimit,
	}
	c.JSON(http.StatusOK, res)
}
