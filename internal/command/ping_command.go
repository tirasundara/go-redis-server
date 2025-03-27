package command

import "github.com/codecrafters-io/redis-starter-go/internal/resp"

// PingCommand implements the PING command
type PingCommand struct{}

// Ensure PingCommand implements Handler interface
var _ Handler = (*PingCommand)(nil)

func (c *PingCommand) Name() string {
	return "PING"
}

func (c *PingCommand) Execute(args []string) resp.RedisValue {
	if len(args) > 1 {
		return resp.Error{Value: "ERR wrong number of arguments for 'ping' command"}
	}

	if len(args) == 1 {
		return resp.BulkString{Value: args[0]}
	}

	return resp.SimpleString{Value: "PONG"}
}
