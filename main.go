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

	var line string

	for {
		data := make([]byte, 8)
		_, err := f.Read(data)
		if err != nil {
			return
		}

		for _, s := range string(data) {
			if s == '\n' {
				fmt.Printf("read: %s\n", line)
				line = ""
				continue
			}
			line += string(s)
		}

	}
}
