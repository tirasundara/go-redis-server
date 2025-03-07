package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

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
