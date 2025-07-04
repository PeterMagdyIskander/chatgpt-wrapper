package models

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// Choice represents a choice in the OpenAI response
type Choice struct {
	Index int `json:"index"`
	Delta struct {
		Content string `json:"content"`
		Role    string `json:"role,omitempty"`
	} `json:"delta"`
	FinishReason *string `json:"finish_reason"`
}

// OpenAIStreamResponse represents a streaming response from OpenAI
type OpenAIStreamResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}
