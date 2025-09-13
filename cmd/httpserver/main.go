package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/samirhembrom/httpfromtcp/internal/request"
	"github.com/samirhembrom/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w io.Writer, req *request.Request) server.HandlerError {
	target := req.RequestLine.RequestTarget
	if target == "/yourproblem" {
		return server.HandlerError{StatusCode: 400, Message: "Your problem is not my problem\n"}
	} else if target == "/myproblem" {
		return server.HandlerError{StatusCode: 500, Message: "Woopsie, my bad\n"}
	}
	w.Write([]byte("All good, frfr\n"))
	return server.HandlerError{}
}
