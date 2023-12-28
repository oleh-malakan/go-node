package main

import (
	"log"
	"net"

	"github.com/oleh-malakan/go-node"
)

func main() {
	server, err := node.Run(0, &net.UDPAddr{})
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	listener, err := server.Listen("Hello, World!")
	if err != nil {
		log.Print(err)

		return
	}
	defer listener.Close()

	stream, err := listener.Accept()
	if err != nil {
		log.Print(err)

		return
	}
	defer stream.Close()

	b, err := stream.Receive()
	if err != nil {
		log.Print(err)

		return
	}

	b = append(b, []byte(", World!")...)

	err = stream.Send(b)
	if err != nil {
		log.Print(err)

		return
	}
}
