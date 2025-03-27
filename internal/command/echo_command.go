package command

import "github.com/codecrafters-io/redis-starter-go/internal/resp"

// EchoCommand implements the ECHO command
type EchoCommand struct{}

// Ensure EchoCommand implements Handler interface
var _ Handler = (*EchoCommand)(nil)

func (c *EchoCommand) Name() string {
	return "ECHO"
}

func (c *EchoCommand) Execute(args []string) resp.RedisValue {
	if len(args) != 1 {
		return resp.Error{Value: "ERR wrong number of arguments for 'echo' command"}
	}

	return resp.BulkString{Value: args[0]}
}
