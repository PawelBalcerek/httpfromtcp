package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/PawelBalcerek/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to start listening port %d: %w", port, err)
	}

	s := &Server{listener: listener}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	if s.closed.CompareAndSwap(false, true) {
		return errors.New("listener already closed")
	}

	if s.listener == nil {
		return errors.New("listener is nil")
	}

	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("an error has occurred while accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	if err := response.WriteStatusLine(conn, response.Ok); err != nil {
		log.Printf("an error has occurred while writing status line: %v", err)
	}

	h := response.GetDefaultHeaders(0)
	if err := response.WriteHeaders(conn, h); err != nil {
		log.Printf("an error has occurred while writing default headers: %v", err)
	}
}
