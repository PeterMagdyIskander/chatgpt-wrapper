package handlers

import (
	"bff/models"
	"bff/services"
	"net/http"

	"github.com/gin-gonic/gin"
)


type KeywordHandlers struct {
	keywordService *services.KeywordService
}


func NewKeywordHandlers(keywordService *services.KeywordService) *KeywordHandlers {
	return &KeywordHandlers{
		keywordService: keywordService,
	}
}


func (h *KeywordHandlers) PostKeywords(c *gin.Context) {
	var req models.KeywordRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON. Expected { \"keywords\": [\"word1\", \"word2\"] }",
		})
		return
	}

	if len(req.Keywords) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The 'keywords' array cannot be empty.",
		})
		return
	}

	h.keywordService.AddWords(req.Keywords)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Lemmatized keywords added",
		"count":   len(req.Keywords),
	})
}


func (h *KeywordHandlers) GetKeywords(c *gin.Context) {
	keywords := h.keywordService.GetAllKeywords()
	c.JSON(http.StatusOK, keywords)
}
