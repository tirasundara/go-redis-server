package resp

import "strconv"

// RedisValue represents a RESP (Redis Serialization Protocol) value
type RedisValue interface {
	// Serialize returns the RESP representation of the value
	Serialize() []byte
}

// SimpleString represents RESP Simple String
type SimpleString struct {
	Value string
}

func (s SimpleString) Serialize() []byte {
	return []byte("+" + s.Value + "\r\n")
}

// Error represents a RESP error
type Error struct {
	Value string
}

// Serialize returns the RESP representation of an Error
func (e Error) Serialize() []byte {
	return []byte("-" + e.Value + "\r\n")
}

// BulkString represents a RESP Bulk String
type BulkString struct {
	Value string
}

// Serialize returns the RESP representation of a Bulk String
func (b BulkString) Serialize() []byte {
	if b.Value == "" {
		return []byte("$-1\r\n")
	}

	return []byte("$" + strconv.Itoa(len(b.Value)) + "\r\n" + b.Value + "\r\n")
}

// Array represents a RESP Array
type Array struct {
	Values []RedisValue
}

// Serialize returns the RESP representation of an Array
func (a Array) Serialize() []byte {
	result := []byte("*" + strconv.Itoa(len(a.Values)) + "\r\n")
	for _, value := range a.Values {
		result = append(result, value.Serialize()...)
	}
	return result
}

// NullBulkString represents a RESP Null Bulk String
var NullBulkString = BulkString{Value: ""}

// CustomResponse allows sending raw RESP data for compatibility with specific implementations
type CustomResponse struct {
	Data []byte
}

// Serialize returns the raw bytes for the custom response
func (c *CustomResponse) Serialize() []byte {
	return c.Data
}
