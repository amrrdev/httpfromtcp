package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	channel := make(chan string)

	go func() {
		defer f.Close()
		defer close(channel)
		currentLine := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Printf("error: %s\n", err.Error())
				}
				break
			}

			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				currentLine += string(data[:i])
				data = data[i+1:]
				channel <- currentLine
				currentLine = ""
			}

			currentLine += string(data)
		}

		if len(currentLine) != 0 {
			channel <- currentLine
		}
	}()

	return channel
}

const filePath = "messages.txt"

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("can not open %s: %s\n", filePath, err)
	}

	channel := getLinesChannel(f)
	for line := range channel {
		fmt.Println("read:", line)
	}

}
