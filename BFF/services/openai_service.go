package services

import (
	"bff/models"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAIService handles communication with OpenAI API
type OpenAIService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewOpenAIService creates a new OpenAI service
func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// StreamCompletion sends a message to OpenAI and streams the response
func (s *OpenAIService) StreamCompletion(userMessage string, responseChan chan<- string, errorChan chan<- error) {
	defer close(responseChan)
	defer close(errorChan)

	// Prepare the request
	requestBody := models.OpenAIRequest{
		Model: "gpt-3.5-turbo", // You can change this to gpt-4 if needed
		Messages: []models.Message{
			{
				Role:    "user",
				Content: userMessage,
			},
		},
		Stream: true,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		errorChan <- fmt.Errorf("failed to marshal request: %w", err)
		return
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", s.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		errorChan <- fmt.Errorf("failed to create request: %w", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		errorChan <- fmt.Errorf("failed to make request: %w", err)
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errorChan <- fmt.Errorf("OpenAI API error: %d - %s", resp.StatusCode, string(body))
		return
	}

	// Read the streaming response
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for data lines
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// Check for end of stream
			if data == "[DONE]" {
				break
			}

			// Parse the JSON response
			var streamResp models.OpenAIStreamResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				// Skip malformed JSON, don't break the stream
				continue
			}

			// Extract content from the response
			if len(streamResp.Choices) > 0 {
				content := streamResp.Choices[0].Delta.Content
				if content != "" {
					responseChan <- content
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		errorChan <- fmt.Errorf("error reading stream: %w", err)
	}
}

// ValidateAPIKey checks if the API key is valid by making a simple request
func (s *OpenAIService) ValidateAPIKey() error {
	req, err := http.NewRequest("GET", s.baseURL+"/models", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("invalid OpenAI API key")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OpenAI API error: %d", resp.StatusCode)
	}

	return nil
}
