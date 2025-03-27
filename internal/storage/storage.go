package storage

// Storage defines the interface for data persistence operations
type Storage interface {
	// Set stores value with no expiration
	Set(key, value string)

	// SetPX stores value with expiration in milliseconds
	SetPX(key, value string, millisecond int)

	// Get retrieves value, returning the value and whether it exists
	Get(key string) (string, bool)

	// GetKeys returns all keys in the storage
	GetKeys() []string

	// Delete removes a key from storage
	Delete(key string) bool

	// LoadRDB loads data from an RDB file
	LoadRDB(filename string) error
}
