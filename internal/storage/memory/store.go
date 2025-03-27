package memory

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/storage"
	"github.com/hdt3213/rdb/model"
	"github.com/hdt3213/rdb/parser"
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

// Ensure Store implements the Storage interface
var _ storage.Storage = (*Store)(nil)

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

// Delete removes a key from the store
func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.data[key]
	if !exists {
		return false
	}

	delete(s.data, key)
	return true
}

// LoadRDB loads data from an RDB file
func (s *Store) LoadRDB(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := parser.NewDecoder(file)

	// Parse RDB file and process entries
	err = decoder.Parse(func(object model.RedisObject) bool {
		key := object.GetKey()
		expiry := object.GetExpiration()

		switch value := object.(type) {
		case *model.StringObject:
			val := string(value.Value)
			if expiry != nil {
				if !time.Now().After(*expiry) { // not expired yet
					expTimeMilli := time.Until(*expiry).Milliseconds()
					s.SetPX(key, val, int(expTimeMilli))
				}
			} else {
				s.Set(key, val)
			}
		default:
			fmt.Printf("Unknown type for key: %s\n", key)
		}

		return true // continue parsing
	})
	if err != nil {
		return fmt.Errorf("Failed to parse RDB: %w", err)
	}

	return nil
}
