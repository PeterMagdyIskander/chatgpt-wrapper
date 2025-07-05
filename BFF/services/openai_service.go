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


type OpenAIService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}


func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}


func (s *OpenAIService) StreamCompletion(userMessage string, responseChan chan<- string, errorChan chan<- error) {
	defer close(responseChan)
	defer close(errorChan)


	requestBody := models.OpenAIRequest{
		Model: "gpt-4o-mini", // You can change this to gpt-4 if needed
		Messages: []models.Message{
			{
				Role: "assistant",
				Content: `Task Instructions:
You will be provided with a question from the user, and you need to provide an answer based on the data provided.
If the question can be answered with the data provided, you should provide a direct answer.
If the question requires reasoning or analysis, you should provide a detailed explanation of your reasoning process and the steps you took to arrive at your answer.
Your answer should always be structured and use fun emojis.`,
			},
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


	req, err := http.NewRequest("POST", s.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		errorChan <- fmt.Errorf("failed to create request: %w", err)
		return
	}


	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)


	resp, err := s.client.Do(req)
	if err != nil {
		errorChan <- fmt.Errorf("failed to make request: %w", err)
		return
	}
	defer resp.Body.Close()


	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errorChan <- fmt.Errorf("OpenAI API error: %d - %s", resp.StatusCode, string(body))
		return
	}


	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}


		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")


			if data == "[DONE]" {
				break
			}


			var streamResp models.OpenAIStreamResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {

				continue
			}


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
