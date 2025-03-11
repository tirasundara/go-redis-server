package memory

import (
	"sync"
	"time"
)

// Store represents an in-memory Redis-like data store
type Store struct {
	mu   sync.RWMutex
	data map[string]entry
}

// entry represents a value in the store
type entry struct {
	value      string
	expiryTime *time.Time
}

// NewStore creates a new in-memory store
func NewStore() *Store {
	return &Store{
		data: make(map[string]entry),
	}
}

// Set sets a key to a string value
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = entry{value: value}
}

// SetPX sets a key with an expiration time in milliseconds
func (s *Store) SetPX(key, value string, millisecond int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	expiryTime := time.Now().Add(time.Duration(millisecond) * time.Millisecond)
	s.data[key] = entry{
		value:      value,
		expiryTime: &expiryTime,
	}
}

// Get retrieves a string value for a key
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, ok := s.data[key]
	if !ok {
		return "", false
	}

	// Check if expired
	if e.expiryTime != nil && time.Now().After(*e.expiryTime) {
		// Should delete, but we're holding a read lock
		// In a more complete implementation, we'd use a background goroutine for cleanup
		return "", false
	}

	// str, ok := e.value.(string)
	str := e.value
	return str, ok
}

func (s *Store) GetKeys() []string {
	keys := make([]string, 0)
	for k := range s.data {
		keys = append(keys, k)
	}

	return keys
}
