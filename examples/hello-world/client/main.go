package main

import (
	"fmt"
	"log"
	"net"

	"github.com/oleh-malakan/go-node"
)

func main() {
	client, err := node.Dial(&net.UDPAddr{})
	if err != nil {
		log.Fatal(err)
	}

	stream, err := client.Connect("Hello, World!")
	if err != nil {
		log.Fatal(err)
	}

	err = stream.Send([]byte("Hello"))
	if err != nil {
		log.Print(err)

		return
	}

	b, err := stream.Receive()
	if err != nil {
		log.Print(err)

		return
	}
	fmt.Println(string(b))

	err = stream.Close()
	if err != nil {
		log.Print(err)

		return
	}
}
