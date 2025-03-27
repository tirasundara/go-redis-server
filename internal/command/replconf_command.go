package command

import "github.com/codecrafters-io/redis-starter-go/internal/resp"

// ReplConfCommand implements the REPLCONF command for replication
type ReplConfCommand struct{}

// Ensure ReplConfCommand implements Handler
var _ Handler = (*ReplConfCommand)(nil)

func NewReplConfCommand() *ReplConfCommand {
	return &ReplConfCommand{}
}

func (c *ReplConfCommand) Name() string {
	return "REPLCONF"
}

func (c *ReplConfCommand) Execute(args []string) resp.RedisValue {
	// For now, simply acknowledge all REPLCONF commands
	return resp.SimpleString{Value: "OK"}
}
