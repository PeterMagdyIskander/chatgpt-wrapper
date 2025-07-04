package main

import (
	"net/http"
	"sync"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/gin-gonic/gin"
)

type userMessages struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	UserId  string `json:"userId"`
}
type MessageStore struct {
	messages []userMessages
	mu       sync.Mutex
}

func (s *MessageStore) AddMessage(msg userMessages) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, msg)
}

func (s *MessageStore) GetAllMessages() []userMessages {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]userMessages{}, s.messages...) // safe copy
}

func postMessages(store *MessageStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newMessage userMessages

		if err := c.BindJSON(&newMessage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		store.AddMessage(newMessage)

		c.IndentedJSON(http.StatusCreated, newMessage)
	}
}

func getMessages(store *MessageStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		messages := store.GetAllMessages()
		c.IndentedJSON(http.StatusOK, messages)
	}
}

// Request struct
type KeywordRequest struct {
	Keywords []string `json:"keywords"`
}

// Set implementation
type Set map[string]struct{}

func (s Set) Add(item string) {
	s[item] = struct{}{}
}

func (s Set) Contains(item string) bool {
	_, ok := s[item]
	return ok
}

func (s Set) Remove(item string) {
	delete(s, item)
}

func (s Set) Size() int {
	return len(s)
}

func (s Set) Values() []string {
	result := make([]string, 0, len(s))
	for k := range s {
		result = append(result, k)
	}
	return result
}

// Thread-safe wrapper for lemmatized keywords
type LemmatizedKeywordStore struct {
	set        Set
	mu         sync.Mutex
	lemmatizer *golem.Lemmatizer
}

// Constructor
func NewLemmatizedKeywordStore() (*LemmatizedKeywordStore, error) {
	lemmatizer, err := golem.New(en.New())
	if err != nil {
		return nil, err
	}

	return &LemmatizedKeywordStore{
		set:        make(Set),
		lemmatizer: lemmatizer,
	}, nil
}

// Lemmatize and store keywords
func (s *LemmatizedKeywordStore) AddWords(words []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, word := range words {
		lemma := s.lemmatizer.Lemma(word)
		s.set.Add(lemma)
	}
}

// Return all lemmatized values
func (s *LemmatizedKeywordStore) GetAll() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.set.Values()
}

// POST handler
func postLemmatizedKeywords(store *LemmatizedKeywordStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req KeywordRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON. Expected { \"keywords\": [\"word1\", \"word2\"] }"})
			return
		}

		if len(req.Keywords) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The 'keywords' array cannot be empty."})
			return
		}

		store.AddWords(req.Keywords)
		c.JSON(http.StatusCreated, gin.H{
			"message": "Lemmatized keywords added",
			"count":   len(req.Keywords),
		})
	}
}

// GET handler
func getLemmatizedKeywords(store *LemmatizedKeywordStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		lemmas := store.GetAll()
		c.JSON(http.StatusOK, lemmas)
	}
}

func main() {
	router := gin.Default()
	store := &MessageStore{}

	lemmatizedStore, err := NewLemmatizedKeywordStore()
	if err != nil {
		panic("failed to initialize lemmatizer: " + err.Error())
	}

	router.POST("/lemmatized-keywords", postLemmatizedKeywords(lemmatizedStore))
	router.GET("/lemmatized-keywords", getLemmatizedKeywords(lemmatizedStore))

	router.POST("/messages", postMessages(store))
	router.GET("/messages", getMessages(store))

	router.Run("localhost:8081")
}
