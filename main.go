package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const filePath = "messages.txt"

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("can not open %s: %s\n", filePath, err)
	}
	defer file.Close()

	currentLine := ""

	for {
		data := make([]byte, 8)
		n, err := file.Read(data)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}

		data = data[:n]
		if i := bytes.IndexByte(data, '\n'); i != -1 {
			currentLine += string(data[:i])
			data = data[i+1:]
			fmt.Printf("read: %s\n", currentLine)
			currentLine = ""
		}

		currentLine += string(data)
	}

	if len(currentLine) != 0 {
		fmt.Printf("read: %s\n", currentLine)
	}

}
