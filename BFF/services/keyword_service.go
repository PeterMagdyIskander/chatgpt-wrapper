package services

import (
	"bff/utils"
	"strings"
	"sync"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
)


type KeywordService struct {
	set        utils.Set
	mu         sync.Mutex
	lemmatizer *golem.Lemmatizer
}


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


func (s *KeywordService) AddWords(words []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, word := range words {

		lemma := s.lemmatizer.Lemma(strings.ToLower(word))
		s.set.Add(lemma)
	}
}


func (s *KeywordService) GetAllKeywords() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.set.Values()
}


func (s *KeywordService) ContainsKeyword(word string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	lemma := s.lemmatizer.Lemma(strings.ToLower(word))
	return s.set.Contains(lemma)
}



func (s *KeywordService) CheckTextForKeywords(text string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	words := strings.Fields(strings.ToLower(text))
	var foundKeywords []string

	for _, word := range words {

		cleanWord := strings.Trim(word, ".,!?;:\"'")
		lemma := s.lemmatizer.Lemma(cleanWord)

		if s.set.Contains(lemma) {
			foundKeywords = append(foundKeywords, lemma)
		}
	}

	return foundKeywords
}
