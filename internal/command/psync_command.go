package command

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/replication"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

// PSyncCommand implements the PSYNC command for replication
type PSyncCommand struct {
	replConfig replication.Config
}

// Ensure PSyncCommand implements Handler
var _ Handler = (*PSyncCommand)(nil)

func NewPSyncCommand(replConfig *replication.Config) *PSyncCommand {
	return &PSyncCommand{replConfig: *replConfig}
}

func (c *PSyncCommand) Name() string {
	return "PSYNC"
}

func (c *PSyncCommand) Execute(args []string) resp.RedisValue {
	if len(args) < 2 {
		return resp.Error{Value: "ERR wrong number of arguments for 'psync' command"}
	}

	// For initial replication, reply with FULLRESYNC
	response := fmt.Sprintf("FULLRESYNC %s %d",
		c.replConfig.MasterReplID,
		c.replConfig.MasterReplOffset)

	return resp.SimpleString{Value: response}
}
