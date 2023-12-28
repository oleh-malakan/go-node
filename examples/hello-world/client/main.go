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
	defer client.Close()

	stream, err := client.Connect("Hello, World!")
	if err != nil {
		log.Print(err)

		return
	}
	defer stream.Close()

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
}
