package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer f.Close()

	for {
		data := make([]byte, 8)
		_, err := f.Read(data)
		if err != nil {
			return
		}
		fmt.Printf("read: %s\n", data)
	}
}
