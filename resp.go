package resp

import (
	"bufio"
	"net"
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

func (r *RespConn) ReadMessage() ([]byte, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
}
