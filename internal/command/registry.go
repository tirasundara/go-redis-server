package command

import "strings"

// DefaultRegistry is the default implementation of the command registry
type DefaultRegistry struct {
	handlers map[string]Handler
}

// Ensure DefaultRegistry implements the Registry interface
var _ Registry = (*DefaultRegistry)(nil)

// NewRegistry creates a new command registry
func NewRegistry() *DefaultRegistry {
	return &DefaultRegistry{
		handlers: make(map[string]Handler),
	}
}

// Register adds a command handler to the registry
func (r *DefaultRegistry) Register(handler Handler) {
	r.handlers[strings.ToUpper(handler.Name())] = handler
}

// Get retrieves a command handler by name
func (r *DefaultRegistry) Get(name string) (Handler, bool) {
	handler, ok := r.handlers[strings.ToUpper(name)]
	return handler, ok
}

// GetAll returns all registered handlers
func (r *DefaultRegistry) GetAll() []Handler {
	handlers := make([]Handler, 0, len(r.handlers))
	for _, h := range r.handlers {
		handlers = append(handlers, h)
	}

	return handlers
}
