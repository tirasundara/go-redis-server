package main

import (
	"flag"
	"fmt"
)

// Config struct to store parsed arguments
type Config struct {
	Dir               string
	DbFileName        string
	Port              uint
	ReplicationConfig *ReplicationConfig
}

type ReplicationConfig struct {
	role             string
	replicaOf        string
	masterReplId     string
	masterReplOffset int
}

func NewReplicationConfig(replicaOf, masterReplId string, masterReplOffset int) *ReplicationConfig {
	replCfg := &ReplicationConfig{
		role:             "master",
		replicaOf:        replicaOf,
		masterReplId:     masterReplId,
		masterReplOffset: masterReplOffset,
	}
	if len(replCfg.replicaOf) > 0 {
		replCfg.role = "slave"
	}

	return replCfg
}

func (c *Config) DbFilePath() string {
	return c.Dir + "/" + c.DbFileName
}

func (c *Config) GetReplicationInfo() string {
	info := "# Replication\n"
	info += fmt.Sprintf("role:%s\nmaster_replid:%s\nmaster_repl_offset:%d\n", c.ReplicationConfig.role, c.ReplicationConfig.masterReplId, c.ReplicationConfig.masterReplOffset)

	return info
}

func loadConfig() *Config {
	// Define flags with default value
	dir := flag.String("dir", "/var/lib/redis", "Directory to store database files")
	dbFilename := flag.String("dbfilename", "dump.rdb", "Database filename")
	port := flag.Uint("port", 6379, "Server port number")
	replicaof := flag.String("replicaof", "", "Replication info. If master it's empty")

	// Parse the command-line arguments
	flag.Parse()

	// TODO: change default values later
	replConfig := NewReplicationConfig(*replicaof, "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb", 0)

	// Store values in Config struct
	return &Config{
		Dir:               *dir,
		DbFileName:        *dbFilename,
		Port:              *port,
		ReplicationConfig: replConfig,
	}
}
