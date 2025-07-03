package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Todo represents a todo item
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TodoRequest represents the request payload for creating/updating todos
type TodoRequest struct {
	Title string `json:"title"`
	Done  *bool  `json:"done,omitempty"`
}

// TodoStore manages todos in memory
type TodoStore struct {
	mu    sync.RWMutex
	todos map[int]*Todo
	nextID int
}

// NewTodoStore creates a new todo store
func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos:  make(map[int]*Todo),
		nextID: 1,
	}
}

// Add creates a new todo
func (ts *TodoStore) Add(title string) *Todo {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	todo := &Todo{
		ID:        ts.nextID,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ts.todos[ts.nextID] = todo
	ts.nextID++

	return todo
}

// Update modifies an existing todo
func (ts *TodoStore) Update(id int, title *string, done *bool) (*Todo, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	todo, exists := ts.todos[id]
	if !exists {
		return nil, fmt.Errorf("todo not found")
	}

	if title != nil {
		todo.Title = *title
	}
	if done != nil {
		todo.Done = *done
	}
	todo.UpdatedAt = time.Now()

	return todo, nil
}

// Delete removes a todo
func (ts *TodoStore) Delete(id int) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.todos[id]; !exists {
		return fmt.Errorf("todo not found")
	}

	delete(ts.todos, id)
	return nil
}

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	store *TodoStore
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(store *TodoStore) *TodoHandler {
	return &TodoHandler{store: store}
}

// AddTodo handles POST /todos
func (h *TodoHandler) AddTodo(w http.ResponseWriter, r *http.Request) {
	var req TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	todo := h.store.Add(req.Title)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// EditTodo handles PUT /todos/{id}
func (h *TodoHandler) EditTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var req TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var title *string
	if req.Title != "" {
		title = &req.Title
	}

	todo, err := h.store.Update(id, title, req.Done)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// DeleteTodo handles DELETE /todos/{id}
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	if err := h.store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	store := NewTodoStore()
	handler := NewTodoHandler(store)

	r := mux.NewRouter()
	
	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/todos", handler.AddTodo).Methods("POST")
	api.HandleFunc("/todos/{id}", handler.EditTodo).Methods("PUT")
	api.HandleFunc("/todos/{id}", handler.DeleteTodo).Methods("DELETE")

	// Middleware for CORS and JSON content type
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	})

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}