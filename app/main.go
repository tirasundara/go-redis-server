package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/command"
	"github.com/codecrafters-io/redis-starter-go/internal/config"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/server"
	"github.com/codecrafters-io/redis-starter-go/internal/storage"
	"github.com/codecrafters-io/redis-starter-go/internal/storage/memory"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// Initialize application configuration
	cfg := config.NewConfig()
	cfg.LoadFromArgs()

	// In-memory key-value store
	var store = memory.NewStore()

	// Load RDB file if exists
	err := loadRDBData(store, cfg.DbFilePath())
	if err != nil {
		fmt.Printf("Warning: Failed to load RDB file: %v\n", err)
	}

	// Set up replication if needed
	if cfg.ReplicationConfig.Role == "slave" {
		_, err := cfg.ReplicationConfig.HandshakeWithMaster()
		if err != nil {
			fmt.Printf("Error connecting to master: %v\n", err)
			os.Exit(1)
		}
	}

	// Set up command registry and register commands
	registry := command.NewRegistry()
	registerCommands(registry, store, cfg)

	// Create and start server
	parser := resp.NewParser()
	redisServer := server.NewServer("0.0.0.0", cfg.Port, registry, parser)

	fmt.Printf("Starting Redis server on port %d\n", cfg.Port)
	err = redisServer.Start()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

// registerCommands registers all supported commands with the registry
func registerCommands(registry command.Registry, store storage.Storage, cfg *config.Config) {
	// Basic commands
	registry.Register(&command.PingCommand{})
	registry.Register(&command.EchoCommand{})
	registry.Register(command.NewGetCommand(store))
	registry.Register(command.NewSetCommand(store))
	registry.Register(command.NewKeysCommand(store))

	// Commands that need configuration
	registry.Register(command.NewInfoCommand(cfg))
	registry.Register(command.NewConfigCommand(cfg))

	// Replication-related commands
	registry.Register(command.NewReplConfCommand())
	registry.Register(command.NewPSyncCommand(cfg.ReplicationConfig))

	// TODO: Add more commands here
}

func loadRDBData(store storage.Storage, filename string) error {
	fmt.Printf("Attempting to load RDB from: %s\n", filename)
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("RDB file does not exist. Starting with empty database...")
			return nil
		}
		return err
	}

	return store.LoadRDB(filename)
}
