package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/samirhembrom/httpfromtcp/internal/request"
	"github.com/samirhembrom/httpfromtcp/internal/response"
)

// Server is an HTTP 1.1 server
type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(*response.Writer, *request.Request)

func Serve(port int, f Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
	}
	go s.listen(f)
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen(f Handler) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn, f)
	}
}

func (s *Server) handle(conn net.Conn, f Handler) {
	defer conn.Close()
	requestData, err := request.RequestFromReader(conn)
	if err != nil {
		writeHandlerError(conn, HandlerError{StatusCode: 500, Message: "Internal server error"})
		return
	}

	buf := bytes.Buffer{}
	writers := response.NewWriter(&buf)
	f(writers, requestData)

	conn.Write(buf.Bytes())
	return
}

func writeHandlerError(w io.Writer, h HandlerError) {
	response.WriteStatusLine(w, response.StatusCode(h.StatusCode))
	headers := response.GetDefaultHeaders(len(h.Message))
	response.WriteHeaders(w, headers)
	w.Write([]byte(h.Message))
}
