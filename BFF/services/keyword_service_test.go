package services

import (
	"bff/utils"
	"testing"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupKeywordService(t *testing.T) *KeywordService {
	lemmatizer, err := golem.New(en.New())
	require.NoError(t, err)

	return &KeywordService{
		set:        utils.NewSet(), // Assuming utils.NewSet() exists
		lemmatizer: lemmatizer,
	}
}

func TestAddWords(t *testing.T) {
	t.Run("should add multiple words and lemmatize them", func(t *testing.T) {
		service := setupKeywordService(t)
		words := []string{"running"}

		service.AddWords(words)

		// Verify lemmatized forms are stored
		assert.True(t, service.set.Contains("run")) // running -> run
	})

	t.Run("should handle empty slice and duplicates", func(t *testing.T) {
		service := setupKeywordService(t)

		// Test empty slice
		service.AddWords([]string{})
		assert.Equal(t, 0, service.set.Size())

		// Test duplicates
		words := []string{"run", "running", "runs"}
		service.AddWords(words)

		// All should lemmatize to "run" - only one entry should exist
		assert.True(t, service.set.Contains("run"))
		assert.Equal(t, 1, service.set.Size())
	})
}

func TestContainsKeyword(t *testing.T) {
	t.Run("should return true for existing lemmatized keywords", func(t *testing.T) {
		service := setupKeywordService(t)
		service.AddWords([]string{"running", "cats"})

		// Test various forms of the same word
		assert.True(t, service.ContainsKeyword("run"))
		assert.True(t, service.ContainsKeyword("running"))
		assert.True(t, service.ContainsKeyword("runs"))
		assert.True(t, service.ContainsKeyword("RUNNING"))
		assert.True(t, service.ContainsKeyword("cat"))
		assert.True(t, service.ContainsKeyword("cats"))
	})

	t.Run("should return false for non-existing keywords", func(t *testing.T) {
		service := setupKeywordService(t)
		service.AddWords([]string{"running"})

		assert.False(t, service.ContainsKeyword("swimming"))
		assert.False(t, service.ContainsKeyword("jump"))
		assert.False(t, service.ContainsKeyword(""))
		assert.False(t, service.ContainsKeyword("walk"))
	})
}

func TestCheckTextForKeywords(t *testing.T) {
	t.Run("should find keywords in message", func(t *testing.T) {
		service := setupKeywordService(t)
		service.AddWords([]string{"running", "cats", "swimming"})

		text := "I love running! My cats are swimming, and dogs are barking."
		foundKeywords := service.CheckTextForKeywords(text)

		// Should find lemmatized forms
		assert.Contains(t, foundKeywords, "run")
		assert.Contains(t, foundKeywords, "cat")
		assert.NotContains(t, foundKeywords, "dog") // not in our keyword set
		assert.Len(t, foundKeywords, 3)
	})

	t.Run("should return empty slice when no keywords found", func(t *testing.T) {
		service := setupKeywordService(t)
		service.AddWords([]string{"running", "swimming"})

		text := "The quick brown fox jumps over the lazy dog."
		foundKeywords := service.CheckTextForKeywords(text)

		assert.Empty(t, foundKeywords)
		assert.Len(t, foundKeywords, 0)

		// Test with empty text
		foundKeywords = service.CheckTextForKeywords("")
		assert.Empty(t, foundKeywords)
	})
}
