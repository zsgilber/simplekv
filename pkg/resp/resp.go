package resp

import (
	"bufio"
	"fmt"
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
		return nil, err
	case BulkString:
		fmt.Println("got to bulk string")
		return r.ReadBulkString()
	}
	return nil, nil
}

func (r *RespConn) ReadBulkString() ([]byte, error) {
	n, err := r.readInt()
	if err != nil {
		return nil, err
	}
	fmt.Println(strconv.Itoa(n))
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
