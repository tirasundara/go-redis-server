package command

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

// ConfigProvider defines an interface for accessing configuration values
type ConfigProvider interface {
	GetString(key string) (string, bool)
}

// ConfigCommand implements the CONFIG command
type ConfigCommand struct {
	config ConfigProvider
}

// NewConfigCommand creates a new CONFIG command handler
func NewConfigCommand(config ConfigProvider) *ConfigCommand {
	return &ConfigCommand{config: config}
}

func (c *ConfigCommand) Name() string {
	return "CONFIG"
}

func (c *ConfigCommand) Execute(args []string) resp.RedisValue {
	if len(args) < 1 {
		return resp.Error{Value: "ERR wrong number of arguments for 'config' command"}
	}

	subcommand := strings.ToUpper(args[0])
	switch subcommand {
	case "GET":
		if len(args) != 2 {
			return resp.Error{Value: "ERR wrong number of arguments for 'config get' command"}
		}
		return c.handleConfigGet(args[1])
	default:
		return resp.Error{Value: fmt.Sprintf("ERR Unknown CONFIG subcommand: %s", args[0])}
	}
}

func (c *ConfigCommand) handleConfigGet(pattern string) resp.RedisValue {
	// If pattern is "*", return all config values
	if pattern == "*" {
		// For this simple implementation, we'll just check a few predefined keys
		keys := []string{"dir", "dbfilename", "port"}
		values := make([]resp.RedisValue, 0, len(keys)*2)

		for _, key := range keys {
			value, found := c.config.GetString(key)
			if found {
				values = append(values, resp.BulkString{Value: key})
				values = append(values, resp.BulkString{Value: value})
			}
		}

		return resp.Array{Values: values}
	}

	// Otherwise, return just the requested value
	value, found := c.config.GetString(strings.ToLower(pattern))
	if !found {
		return resp.Array{Values: []resp.RedisValue{}}
	}

	return resp.Array{Values: []resp.RedisValue{
		resp.BulkString{Value: pattern},
		resp.BulkString{Value: value},
	}}
}
