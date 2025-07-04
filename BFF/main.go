package main

import (
	"bff/handlers"
	"bff/services"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Get OpenAI API key from environment variable
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if openaiAPIKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Initialize services
	messageService := services.NewMessageService()

	keywordService, err := services.NewKeywordService()
	if err != nil {
		log.Fatal("Failed to initialize keyword service:", err)
	}

	openaiService := services.NewOpenAIService(openaiAPIKey)

	// Validate OpenAI API key
	if err := openaiService.ValidateAPIKey(); err != nil {
		log.Fatal("Failed to validate OpenAI API key:", err)
	}

	// Initialize handlers
	messageHandlers := handlers.NewMessageHandlers(messageService, keywordService)
	keywordHandlers := handlers.NewKeywordHandlers(keywordService)
	sseHandlers := handlers.NewSSEHandlers(messageService, openaiService)

	// Setup router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Message routes
	router.POST("/messages", messageHandlers.PostMessage)
	router.GET("/messages", messageHandlers.GetMessages)

	// Keyword routes
	router.POST("/lemmatized-keywords", keywordHandlers.PostKeywords)
	router.GET("/lemmatized-keywords", keywordHandlers.GetKeywords)

	// SSE/Streaming routes
	router.GET("/ask-chatgpt", sseHandlers.StreamCompletion)

	// Start server
	log.Println("Server starting on :8081")
	log.Println("Make sure to set OPENAI_API_KEY environment variable")
	if err := router.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
