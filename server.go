package resp

import (
	"net"
	"sync"
)

// Server defines a server for listening to client requests.
type Server struct {
	mu       sync.RWMutex
	handlers map[string]func(resp *RespConn)
}

// NewServer creates a new server with an empty handlers map.
func NewServer() *Server {
	return &Server{
		handlers: make(map[string]func(resp *RespConn)),
	}
}

func (s *Server) AddHandleFunc(command string, handler func(resp *RespConn)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[command] = handler
}

func (s *Server) ListenAndServe(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			//TODO: log error
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) error {
	resp := NewRespConn(conn)

	for {
		resp.ReadByte()
	}
}
