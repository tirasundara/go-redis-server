package command

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/storage"
)

// KeysCommand implements the KEYS command
type KeysCommand struct {
	store storage.Storage
}

// Ensure KeysCommands implements Handler
var _ (Handler) = (*KeysCommand)(nil)

func NewKeysCommand(store storage.Storage) *KeysCommand {
	return &KeysCommand{store: store}
}

func (c *KeysCommand) Name() string {
	return "KEYS"
}

func (c *KeysCommand) Execute(args []string) resp.RedisValue {
	if len(args) != 1 {
		return resp.Error{Value: "ERR wrong number of arguments for 'keys' command"}
	}

	pattern := args[0]
	if pattern != "*" {
		// For simplicity, we only support the "*" pattern for now
		return resp.Error{Value: "ERR pattern not supported"}
	}

	keys := c.store.GetKeys()
	values := make([]resp.RedisValue, len(keys))
	for i, key := range keys {
		values[i] = resp.BulkString{Value: key}
	}

	return resp.Array{Values: values}
}
