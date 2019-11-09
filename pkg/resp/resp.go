package resp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	SimpleString = '+'
	BulkString   = '$'
	Integer      = ':'
	Array        = '*'
	Error        = '-'
)

// TODO: figure out whether i want object or not
// type Object struct {
// 	t       Type
// 	val     []byte
// 	integer int
// 	string  []byte
// 	array   []Object
// }

type respError struct {
	err error
}

func (e respError) Error() string {
	return fmt.Sprintf("resp protocol error: %v", e.err)
}

// RespConn reads and writes RESP values over a connection.
type RespConn struct {
	*bufio.Reader
	*bufio.Writer
	conn net.Conn
}

// NewRespConn returns creates a new RespConn.
func NewRespConn(conn net.Conn) *RespConn {
	return &RespConn{
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
		conn:   conn,
	}
}

func (r *RespConn) ReadValue() ([]byte, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	switch b {
	case Array:
		fmt.Println("got to array")
		return r.ReadArray()
	case BulkString:
		fmt.Println("got to bulk string")
		return r.ReadBulkString()
	}
	return nil, nil
}

func (r *RespConn) ReadBulkString() ([]byte, error) {
	bulkLength, err := r.readInt()
	if err != nil {
		return nil, err
	}
	if bulkLength == -1 {
		return nil, nil // return nil to reprsent NULL BulkString value.
	}

	buf := make([]byte, bulkLength+2)
	_, err = io.ReadFull(r, buf) // read number of bytes as indicated by bulkLength, plus the two line termination bytes.
	if err != nil {
		return nil, err
	}
	if buf[bulkLength] != '\r' || buf[bulkLength+1] != '\n' {
		return nil, respError{errors.New("invalid line termination")}
	}
	return buf, nil
}

func (r *RespConn) ReadArray() ([][]byte, error) {
	arrayLength, err := r.readInt()
	if err != nil {
		return nil, err
	}
	array := make([][]byte, arrayLength)
	for i := 0; i < arrayLength; i++ {
		val, err := r.ReadValue()
		if err != nil {
			return nil, err
		}
		array[i] = val
	}
	return nil, nil
}

func (r *RespConn) readLine() ([]byte, error) {
	b, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return b[:len(b)-2], nil
}

func (r *RespConn) readInt() (int, error) {
	line, err := r.readLine()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

// func (r *RespConn) readLine() ([]byte, error) {
// 	line, err := r.ReadBytes('\n')
// 	if err != nil {
// 		return nil, err
// 	}

// }
