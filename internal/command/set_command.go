package command

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/storage"
)

// SetCommand implements SET command
type SetCommand struct {
	store storage.Storage
}

// Ensure SetCommand implement Handler
var _ Handler = (*SetCommand)(nil)

func NewSetCommand(store storage.Storage) *SetCommand {
	return &SetCommand{store: store}
}

func (c *SetCommand) Name() string {
	return "SET"
}

func (c *SetCommand) Execute(args []string) resp.RedisValue {
	if len(args) < 2 {
		return resp.Error{Value: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0]
	value := args[1]

	// Check for additional options
	if len(args) > 2 {
		// Handle PX option (expiry in milliseconds)
		if len(args) >= 4 && strings.ToUpper(args[2]) == "PX" {
			ms, err := strconv.Atoi(args[3])
			if err != nil {
				return resp.Error{Value: "ERR invalid expire time in 'set' command"}
			}

			c.store.SetPX(key, value, ms)
			return resp.SimpleString{Value: "OK"}
		}

		// We could add support for other options like EX, NX, XX here
		return resp.Error{Value: "ERR syntax error"}
	}

	c.store.Set(key, value)
	return resp.SimpleString{Value: "OK"}
}
