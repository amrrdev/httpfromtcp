package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpAdd, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("can not resolve udp address: %s", err.Error())
	}

	conn, err := net.DialUDP(udpAdd.Network(), nil, udpAdd)
	if err != nil {
		log.Fatal("could not open a connection")
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = conn.Write([]byte(str))
		if err != nil {
			fmt.Println(err)
		}

	}

}
