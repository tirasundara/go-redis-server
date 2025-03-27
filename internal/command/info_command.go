package command

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

// InfoCommand implements the INFO command
type InfoCommand struct {
	config ReplicationInfo
}

// ReplicationInfo provides replication information for INFO command
type ReplicationInfo interface {
	GetReplicationInfo() string
}

// Ensure InfoCommand implements Handler
var _ Handler = (*InfoCommand)(nil)

func NewInfoCommand(config ReplicationInfo) *InfoCommand {
	return &InfoCommand{config: config}
}

func (c *InfoCommand) Name() string {
	return "INFO"
}

func (c *InfoCommand) Execute(args []string) resp.RedisValue {
	// For now, only handle replication info
	section := ""
	if len(args) > 0 {
		section = strings.ToLower(args[0])
	}

	var info string
	if section == "" || section == "replication" {
		info = c.config.GetReplicationInfo()
	} else {
		info = fmt.Sprintf("# %s\r\n", strings.Title(section))
	}

	return resp.BulkString{Value: info}
}
