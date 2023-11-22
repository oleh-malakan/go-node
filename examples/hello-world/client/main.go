package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/oleh-malakan/go-node"
)

func main() {
	client, err := node.Dial(&tls.Config{}, &net.UDPAddr{})
	if err != nil {
		log.Fatal(err)
	}

	connection, err := client.Connect("Hello, World!")
	if err != nil {
		log.Print(err)

		return
	}

	if err := connection.Send([]byte("Hello")); err != nil {
		log.Print(err)

		return
	}

	b, err := connection.Receive()
	if err != nil {
		log.Print(err)

		return
	}
	fmt.Println(string(b))

	if err := connection.Close(); err != nil {
		log.Print(err)

		return
	}
}
