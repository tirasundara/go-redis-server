package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Parser defines the interface for parsing RESP protocol
type Parser interface {
	// ParseCommand parses a RESP command from a reader
	ParseCommand(reader *bufio.Reader) ([]string, error)
}

// DefaultParser is the default implementation of the RESP protocol parser
type DefaultParser struct{}

// NewParser creates a new RESP parser
func NewParser() *DefaultParser {
	return &DefaultParser{}
}

// ParseCommand parses a RESP command from a reader
func (p *DefaultParser) ParseCommand(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "*") {
		// TODO: Handle raw text commands later
		return nil, fmt.Errorf("Invalid RESP format: expected array prefix '*'")
	}

	count, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, fmt.Errorf("Invalid RESP array length: %w", err)
	}

	commands := make([]string, count)
	for i := range count {
		// Read bulk string header (e.g: `$3`)
		line, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "$") {
			return nil, fmt.Errorf("Invalid RESP format: expected bulk string prefix '$'")
		}

		// Parse bulk string length
		length, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid RESP bulk string length: %w", err)
		}

		// Read actual command content
		if length > 0 {
			// Read the exact number of bytes, plus the trailing \r\n
			data := make([]byte, length+2)
			_, err := io.ReadFull(reader, data)
			if err != nil {
				return nil, err
			}
			commands[i] = string(data[:length])

		} else {
			// Empty string case, just read \r\n
			_, err = reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			commands[i] = ""
		}
	}

	return commands, nil
}
