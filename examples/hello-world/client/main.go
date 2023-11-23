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
		log.Fatal(err)
	}

	err = connection.Send([]byte("Hello"))
	if err != nil {
		log.Print(err)

		return
	}

	b, err := connection.Receive()
	if err != nil {
		log.Print(err)

		return
	}
	fmt.Println(string(b))

	err = connection.Close()
	if err != nil {
		log.Print(err)

		return
	}
}
