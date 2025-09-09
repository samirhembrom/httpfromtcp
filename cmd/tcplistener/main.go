package main

import (
	"fmt"
	"log"
	"net"

	"github.com/samirhembrom/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Issue on port: %s %s", port, err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Connection error: %s", err)
		}
		fmt.Printf("Connection has been accepted\n")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("err")
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %v: %v\n", k, v)
		}

		fmt.Printf("Connection has been closed\n")
	}
}
