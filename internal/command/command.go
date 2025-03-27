package command

import "github.com/codecrafters-io/redis-starter-go/internal/resp"

// Handler defines the interface for command handling
type Handler interface {
	// Name returns the command name (e.g., "GET", "SET")
	Name() string

	// Execute runs the command with the given arguments
	Execute(args []string) resp.RedisValue
}

// Registry maintains a mapping of command names to their handlers
type Registry interface {
	// Register adds a command handler to the registry
	Register(handler Handler)

	// Get retrieves a command handler by name
	Get(name string) (Handler, bool)

	// GetAll returns all registered handlers
	GetAll() []Handler
}
