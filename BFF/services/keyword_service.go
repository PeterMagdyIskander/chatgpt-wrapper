package services

import (
	"bff/utils"
	"strings"
	"sync"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
)

// KeywordService handles keyword lemmatization and storage
type KeywordService struct {
	set        utils.Set
	mu         sync.Mutex
	lemmatizer *golem.Lemmatizer
}

// NewKeywordService creates a new keyword service with lemmatizer
func NewKeywordService() (*KeywordService, error) {
	lemmatizer, err := golem.New(en.New())
	if err != nil {
		return nil, err
	}

	return &KeywordService{
		set:        utils.NewSet(),
		lemmatizer: lemmatizer,
	}, nil
}

// AddWords lemmatizes and stores keywords
func (s *KeywordService) AddWords(words []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, word := range words {
		// Convert to lowercase and lemmatize
		lemma := s.lemmatizer.Lemma(strings.ToLower(word))
		s.set.Add(lemma)
	}
}

// GetAllKeywords returns all stored lemmatized keywords
func (s *KeywordService) GetAllKeywords() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.set.Values()
}

// ContainsKeyword checks if a lemmatized keyword exists
func (s *KeywordService) ContainsKeyword(word string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	lemma := s.lemmatizer.Lemma(strings.ToLower(word))
	return s.set.Contains(lemma)
}

// CheckTextForKeywords checks if a text contains any of the stored keywords
// This will be useful for your future message filtering functionality
func (s *KeywordService) CheckTextForKeywords(text string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	words := strings.Fields(strings.ToLower(text))
	var foundKeywords []string

	for _, word := range words {
		// Clean punctuation from word
		cleanWord := strings.Trim(word, ".,!?;:\"'")
		lemma := s.lemmatizer.Lemma(cleanWord)

		if s.set.Contains(lemma) {
			foundKeywords = append(foundKeywords, lemma)
		}
	}

	return foundKeywords
}
