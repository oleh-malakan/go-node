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

	message := string(b)
	if message == "Hello" {
		if err := stream.Send([]byte(message + ", World!")); err != nil {
			log.Print(err)

			return
		}
	}
}
