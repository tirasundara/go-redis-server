package replication

import (
	"fmt"
	"net"
	"time"
)

// Config holds replication configuration
type Config struct {
	Role             string
	ReplicaOf        string
	MasterHost       string
	MasterPort       int
	MasterReplID     string
	MasterReplOffset int
	Port             int
}

// NewConfig creates a new replication configuration
func NewConfig() *Config {
	return &Config{
		Role:             "master",
		MasterReplID:     "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		MasterReplOffset: 0,
	}
}

// SetupReplication configures this instance as a replica of the specified master
func (c *Config) SetupReplication(masterHost string, masterPort int, myPort int) {
	c.Role = "slave"
	c.MasterHost = masterHost
	c.MasterPort = masterPort
	c.Port = myPort
	c.ReplicaOf = fmt.Sprintf("%s %d", masterHost, masterPort)
}

// GetReplicationInfo returns replication information for the INFO command
func (c *Config) GetReplicationInfo() string {
	info := "# Replication\n"
	info += fmt.Sprintf("role:%s\n", c.Role)
	info += fmt.Sprintf("master_replid:%s\n", c.MasterReplID)
	info += fmt.Sprintf("master_repl_offset:%d\n", c.MasterReplOffset)

	if c.Role == "slave" {
		info += fmt.Sprintf("master_host:%s\n", c.MasterHost)
		info += fmt.Sprintf("master_port:%d\n", c.MasterPort)
	}

	return info
}

// HandshakeWithMaster performs the initial replication handshake with the master
func (c *Config) HandshakeWithMaster() (net.Conn, error) {
	if c.Role != "slave" {
		return nil, fmt.Errorf("not configured as a replica")
	}

	addr := net.JoinHostPort(c.MasterHost, fmt.Sprintf("%d", c.MasterPort))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master at %s: %w", addr, err)
	}

	// Send PING
	conn.Write([]byte("*1\r\n$4\r\nping\r\n"))
	time.Sleep(50 * time.Millisecond)

	// Send REPLCONF listening-port
	conn.Write([]byte(fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$%d\r\n%d\r\n",
		len(fmt.Sprintf("%d", c.Port)), c.Port)))
	time.Sleep(50 * time.Millisecond)

	// Send REPLCONF capa psync2
	conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"))
	time.Sleep(50 * time.Millisecond)

	// Send PSYNC
	conn.Write([]byte("*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"))

	return conn, nil
}
