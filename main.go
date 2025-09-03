package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
		linesChan := getLinesChannel(conn)
		for line := range linesChan {
			fmt.Println(line)
		}
		fmt.Printf("Connection has been closed\n")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	line := make(chan string)

	go func() {
		defer f.Close()
		defer close(line)
		currentLineContents := ""
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					line <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := range len(parts) - 1 {
				line <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return line
}
