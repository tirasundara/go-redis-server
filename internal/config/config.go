package config

import (
	"flag"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/replication"
)

// Config represents the application configuration
type Config struct {
	Dir               string
	DbFileName        string
	Port              int
	ReplicationConfig *replication.Config
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		Dir:               "/var/lib/redis",
		DbFileName:        "dump.rdb",
		Port:              6379,
		ReplicationConfig: replication.NewConfig(),
	}
}

// DbFilePath returns the full path to the RDB file
func (c *Config) DbFilePath() string {
	return filepath.Join(c.Dir, c.DbFileName)
}

// LoadFromArgs loads configuration from command-line arguments
func (c *Config) LoadFromArgs() {
	// Define flags with default values
	dir := flag.String("dir", c.Dir, "Directory to store database files")
	dbFilename := flag.String("dbfilename", c.DbFileName, "Database filename")
	port := flag.Int("port", c.Port, "Server port number")
	replicaOf := flag.String("replicaof", "", "Master host and port for replication (e.g., '127.0.0.1 6379')")

	// Parse the command-line arguments
	flag.Parse()

	// Update config from flags
	c.Dir = *dir
	c.DbFileName = *dbFilename
	c.Port = *port

	// Handle replication configuration
	if *replicaOf != "" {
		fields := strings.Fields(*replicaOf)
		if len(fields) == 2 {
			masterHost := fields[0]
			masterPort, err := strconv.Atoi(fields[1])
			if err == nil {
				c.ReplicationConfig.SetupReplication(masterHost, masterPort, c.Port)
			}
		}
	}
}

// GetString returns a configuration value as a string
func (c *Config) GetString(key string) (string, bool) {
	switch key {
	case "dir":
		return c.Dir, true
	case "dbfilename":
		return c.DbFileName, true
	case "port":
		return strconv.Itoa(c.Port), true
	default:
		return "", false
	}
}

// GetReplicationInfo returns the replication information
func (c *Config) GetReplicationInfo() string {
	return c.ReplicationConfig.GetReplicationInfo()
}
