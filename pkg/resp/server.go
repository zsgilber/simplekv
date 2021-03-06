package resp

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/zsgilber/simplekv/pkg/kv"
)

// Server defines a server for listening to client requests.
type Server struct {
	mu       sync.RWMutex
	handlers map[string]func(resp *RespConn)
	store    kv.Store
}

// NewServer creates a new server with an empty handlers map.
func NewServer(store kv.Store) *Server {
	return &Server{
		handlers: make(map[string]func(resp *RespConn)),
		store:    store,
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
			fmt.Println(err)
			continue
		}
		go s.handleConnection(conn, s.store)
	}
}

//TODO: replace the error return with a WriteError for the respConn
func (s *Server) handleConnection(conn net.Conn, store kv.Store) error {
	respConn := NewRespConn(conn)
	command, err := respConn.ReadCommand()
	if err != nil {
		fmt.Println(err)
		return err
	}
	commandName := string(command.Args[0].str)
	fmt.Println(len(commandName))
	switch strings.ToLower(commandName) {
	case "set":
		if err := store.Set(string(command.Args[1].str), string(command.Args[2].str)); err != nil {
			return err
		}
		respConn.WriteSimpleString("OK")
	case "get":
		value, err := store.Get(string(command.Args[1].str))
		fmt.Println(value)
		if err != nil {
			return err
		}
	case "ping":
		respConn.WriteSimpleString("PONG")
	default:
		fmt.Println("unknown command")
	}
	return nil
}
