package command

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/storage"
)

// GetCommand implements the GET command
type GetCommand struct {
	store storage.Storage
}

func NewGetCommand(store storage.Storage) *GetCommand {
	return &GetCommand{store: store}
}

func (c *GetCommand) Name() string {
	return "GET"
}

func (c *GetCommand) Execute(args []string) resp.RedisValue {
	if len(args) != 1 {
		return resp.Error{Value: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0]
	value, exists := c.store.Get(key)
	if !exists {
		// The original code returns "$-1\r\n" directly
		// For compatibility with the CodeCrafters test, we need to use a custom implementation
		// This creates a custom response that matches exactly what the original code returns
		return &resp.CustomResponse{Data: []byte("$-1\r\n")}
	}

	return resp.BulkString{Value: value}
}
