package resp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

type Type byte

const (
	SimpleString Type = '+'
	BulkString   Type = '$'
	Integer      Type = ':'
	Array        Type = '*'
	Error        Type = '-'
)

// TODO: figure out whether i want object or not
type Object struct {
	t       Type
	integer int
	str     []byte
	array   []Object
	null    bool
}

var (
	nullObject = Object{null: true}
)

type Command struct {
	Args []Object
}

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

func (r *RespConn) ReadCommand() (Command, error) {
	object, err := r.ReadObject()
	if err != nil {
		return Command{}, err
	}
	return Command{
		Args: object.array,
	}, nil
}

func (r *RespConn) ReadObject() (Object, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nullObject, err
	}
	switch b {
	case '*':
		return r.ReadArray()
	case '$':
		return r.ReadBulkString()
	}
	return nullObject, nil
}

func (r *RespConn) ReadBulkString() (Object, error) {
	bulkLength, err := r.readInt()
	if err != nil {
		return nullObject, err
	}
	if bulkLength == -1 {
		return nullObject, nil
	}

	buf := make([]byte, bulkLength+2)
	_, err = io.ReadFull(r, buf) // read number of bytes as indicated by bulkLength, plus the two line termination bytes.
	if err != nil {
		return nullObject, err
	}
	if buf[bulkLength] != '\r' || buf[bulkLength+1] != '\n' {
		return nullObject, respError{errors.New("invalid line termination")}
	}
	return Object{
		t:   BulkString,
		str: buf[:bulkLength],
	}, nil
}

func (r *RespConn) ReadArray() (Object, error) {
	arrayLength, err := r.readInt()
	if err != nil {
		return nullObject, err
	}
	array := make([]Object, arrayLength)
	for i := 0; i < arrayLength; i++ {
		val, err := r.ReadObject()
		if err != nil {
			return nullObject, err
		}
		array[i] = val
	}
	return Object{
		t:     Array,
		array: array,
	}, nil
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
