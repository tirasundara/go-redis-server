package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hdt3213/rdb/model"
	"github.com/hdt3213/rdb/parser"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

type Entry struct {
	Value      string
	ExpiryTime *time.Time
}

// In-memory key-value store
// TODO: add mutex later
var store = make(map[string]Entry)

var config *Config

func main() {
	// Load app configuration
	config = loadConfig()

	// Load entries from RDB file (if exists)
	err := loadRDB(config.DbFilePath())
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Failed to load RDB:", err)
			os.Exit(1)
		}
		fmt.Println("File does not exist. Start from fresh instead")
	}

	// Start the server
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.Port))
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Redis server running on port %d...\n", config.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func loadRDB(filename string) error {
	// writeRDB()
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := parser.NewDecoder(file)

	// Parse RDB file and process entries
	err = decoder.Parse(func(object model.RedisObject) bool {
		key := object.GetKey()
		expiry := object.GetExpiration()

		switch value := object.(type) {
		case *model.StringObject:
			entry := Entry{Value: string(value.Value)}
			if expiry != nil {
				entry.ExpiryTime = expiry
			}
			store[key] = entry
		default:
			fmt.Printf("Unknown type for key: %s\n", key)
		}

		return true // continue parsing
	})

	if err != nil {
		return fmt.Errorf("failed to parse RDB: %w", err)
	}

	return nil
}

// parseRESP reads the incoming data from the client and parses it into a RESP array.
// It returns the parsed RESP array or an error if the parsing fails.
func parseRESP(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	fmt.Println("line:", line)
	if !strings.HasPrefix(line, "*") {
		// TODO: Handle raw text commands later
		return nil, fmt.Errorf("invalid RESP format")
	}

	count, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, err
	}

	commands := make([]string, count)
	for i := 0; i < count; i++ {
		// skip bulk string header (e.g: `$3`)
		_, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		// read actual command
		line, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		commands[i] = line
	}

	return commands, nil
}

func executeCommand(commands []string) string {
	if len(commands) == 0 {
		return "-ERR unknown command"
	}

	switch strings.ToUpper(commands[0]) {
	case "PING":
		return "+PONG\r\n"

	case "ECHO":
		if len(commands) < 2 {
			return "-ERR wrong number of arguments\r\n"
		}
		return fmt.Sprintf("$%d\r\n%s\r\n", len(commands[1]), commands[1])

	case "SET":
		if len(commands) < 3 {
			return "-ERR wrong number of arguments for 'SET'\r\n"
		}

		key := commands[1]
		value := commands[2]
		entry := Entry{Value: value}
		if len(commands) == 5 && strings.ToUpper(commands[3]) == "PX" {
			ms, err := strconv.Atoi(commands[4])
			if err != nil {
				return "-ERR invalid expiry duration value\r\n"
			}
			expiry := time.Now().Add(time.Duration(ms) * time.Millisecond)
			entry.ExpiryTime = &expiry
		}

		// TODO: implement lock & unlock
		store[key] = entry

		return "+OK\r\n"

	case "GET":
		if len(commands) < 2 {
			return "-ERR wrong number of arguments for 'GET'\r\n"
		}

		key := commands[1]
		entry, found := store[key]
		if !found {
			return "$-1\r\n" // Null response if key doesn't exist
		}

		if entry.ExpiryTime != nil && time.Now().After(*entry.ExpiryTime) {
			delete(store, key) // Delete expired key-value
			return "$-1\r\n"
		}
		return fmt.Sprintf("$%d\r\n%s\r\n", len(entry.Value), entry.Value)

	case "CONFIG":
		n := len(commands)
		if n != 3 {
			return "-ERR wrong number of arguments for 'CONFIG'\r\n"
		}

		if strings.ToUpper(commands[1]) != "GET" {
			return fmt.Sprintf("-ERR Unknown CONFIG command: %s \r\n", commands[1])
		}

		return handleConfigGet(commands[2])

	case "KEYS":
		if len(commands) < 2 {
			return "-ERR wrong number of arguments for 'KEYS'\r\n"
		}

		if commands[1] != "*" {
			return "-ERR invalid argument for 'KEYS'\r\n"
		}

		var response string
		response = fmt.Sprintf("*%d\r\n", len(store))
		for k := range store {
			response += fmt.Sprintf("$%d\r\n%s\r\n", len(k), k)
		}
		return response

	case "INFO":
		info := "# Replication\nrole:"
		role := "master"
		if config.ReplicaOf != "" {
			role = "slave"
		}
		info += role
		return fmt.Sprintf("$%d\r\n%s\r\n", len(info), info)
	default:
		return "+PONG\r\n" // TODO: may change later
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		command, err := parseRESP(reader)
		if err != nil {
			if err == io.EOF {
				conn.Close()
				fmt.Println("EOF detected. Connection closed")
				return
			}

			fmt.Println("Error parsing request: ", err.Error())
			return
		}

		fmt.Println("Received command:", command)
		response := executeCommand(command)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			return
		}
		fmt.Printf("Sent: %s\n", response)
	}
}

func handleConfigGet(key string) string {
	var response string
	var configMap = map[string]string{
		"dir":        config.Dir,
		"dbFilename": config.DbFileName,
		"port":       fmt.Sprintf("%d", config.Port),
	}
	if key == "*" { // Return all config values
		response = fmt.Sprintf("*%d\r\n", len(configMap)*2)
		for k, v := range configMap {
			response += fmt.Sprintf("$%d\r\n%s\r\n$%d\r\n%s\r\n", len(k), k, len(v), v)
		}
		return response
	}

	value, found := configMap[key]
	if !found {
		return "$-1\r\n" // RESP Null Bulk String (key not found)
	}
	response = fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(key), key, len(value), value)

	return response
}
