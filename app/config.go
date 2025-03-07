package main

import "flag"

// Config struct to store parsed arguments
type Config struct {
	Dir        string
	DbFileName string
	Port       uint
	ReplicaOf  string
}

func (c *Config) DbFilePath() string {
	return c.Dir + "/" + c.DbFileName
}

func loadConfig() *Config {
	// Define flags with default value
	dir := flag.String("dir", "/var/lib/redis", "Directory to store database files")
	dbFilename := flag.String("dbfilename", "dump.rdb", "Database filename")
	port := flag.Uint("port", 6379, "Server port number")
	replicaof := flag.String("replicaof", "", "Replication info. If master it's empty")

	// Parse the command-line arguments
	flag.Parse()

	// Store values in Config struct
	return &Config{
		Dir:        *dir,
		DbFileName: *dbFilename,
		Port:       *port,
		ReplicaOf:  *replicaof,
	}
}
