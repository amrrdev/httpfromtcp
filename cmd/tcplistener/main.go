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

		fmt.Println("Request Line: ", r.RequestLine)
		fmt.Println("Request Method: ", r.RequestLine.HttpMethod)
		fmt.Println("Request Resource Path: ", r.RequestLine.RequestTarget)
		fmt.Println("Request Version: ", r.RequestLine.HttpVersion)

	}

}
