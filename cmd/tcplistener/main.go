package main

import (
	"fmt"
	"log"
	"net"

	"github.com/amrrdev/httpfromtcp/internal/request"
)

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

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Request Line: ")
		fmt.Println("- Method: ", r.RequestLine.HttpMethod)
		fmt.Println("- Resource Path: ", r.RequestLine.RequestTarget)
		fmt.Println("- Version: ", r.RequestLine.HttpVersion)

		fmt.Println("Headers: ")
		for key, value := range r.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
	}

}
