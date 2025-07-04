package models

// KeywordRequest represents the request structure for adding keywords
type KeywordRequest struct {
	Keywords []string `json:"keywords"`
}
