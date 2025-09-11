package server

import (
	"fmt"
	"net"
)

type Server struct {
	l net.Listener
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return &Server{}, err
	}
	server := &Server{l: l}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	return s.l.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.l.Accept()
		if err != nil {
			return
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	resp := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!\n"
	conn.Write([]byte(resp))
}
