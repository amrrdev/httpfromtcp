package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

func getLinesChannel(conn net.Conn) <-chan string {
	channel := make(chan string)

	go func() {
		defer conn.Close()
		defer close(channel)
		currentLine := ""
		for {
			data := make([]byte, 8)
			n, err := conn.Read(data)
			if err != nil {
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

func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalln("can not start listing on tcp connection", "error", err)
	}

	fmt.Println("start accepting connections...")

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("can not accept connection")
			continue
		}

		channel := getLinesChannel(conn)
		for line := range channel {
			fmt.Println("read:", line)
		}

	}

}
