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
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

type Entry struct {
	Value      string
	ExpiryTime time.Time // Zero value means no expiration
}

// In-memory key-value store
// TODO: add mutex later
var store = make(map[string]Entry)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Redis server running on port 6379...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
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

		var expiry time.Time
		key := commands[1]
		value := commands[2]
		if len(commands) == 5 && strings.ToUpper(commands[3]) == "PX" {
			ms, err := strconv.Atoi(commands[4])
			if err != nil {
				return "-ERR invalid expiry duration value\r\n"
			}
			expiry = time.Now().Add(time.Duration(ms) * time.Millisecond)
		}
		// TODO: implement lock & unlock
		entry := Entry{Value: value, ExpiryTime: expiry}
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

		if !entry.ExpiryTime.IsZero() && time.Now().After(entry.ExpiryTime) {
			delete(store, key) // Delete expired key-value
			return "$-1\r\n"
		}
		return fmt.Sprintf("$%d\r\n%s\r\n", len(entry.Value), entry.Value)

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
