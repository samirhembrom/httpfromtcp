package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const port = "localhost:42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatalf("Error resolving on port %s message: %s", port, udpAddr)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Error dialing on port %s message: %s", port, err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Issue reading input %s\n", err)
		}

		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Fatalf("error printing output %s\n", err)
		}

	}
}
