package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/command"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

// Server represents a Redis server
type Server struct {
	host     string
	port     int
	commands command.Registry
	parser   resp.Parser
}

// NewServer creates a new Redis server
func NewServer(host string, port int, commands command.Registry, parser resp.Parser) *Server {
	return &Server{
		host:     host,
		port:     port,
		commands: commands,
		parser:   parser,
	}
}

// Start starts the Redis server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind to %s: %w", addr, err)
	}
	defer listener.Close()

	fmt.Printf("Redis server running on %s...\n", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}

		go s.handleConnection(conn)
	}
}

// handleConnection processes client connections
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// Parse incoming command
		args, err := s.parser.ParseCommand(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				return
			}
			fmt.Printf("Error parsing command: %v\n", err)
			return
		}

		if len(args) == 0 {
			continue
		}

		// Find command handler
		handlerName := strings.ToUpper(args[0])
		handler, found := s.commands.Get(handlerName)

		var response resp.RedisValue
		if !found {
			response = resp.Error{Value: fmt.Sprintf("ERR unknown command '%s'", handlerName)}
		} else {
			// Execute command with arguments (skip the command name)
			response = handler.Execute(args[1:])
		}

		// Send response
		_, err = conn.Write(response.Serialize())
		if err != nil {
			fmt.Printf("Error writing response: %v\n", err)
			return
		}
	}
}
