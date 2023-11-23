package main

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/oleh-malakan/go-node"
)

func main() {
	server, err := node.New(&tls.Config{}, &net.UDPAddr{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = server.Handler("Hello, World!", func(connection *node.Connection) {
		b, err := connection.Receive()
		if err != nil {
			log.Print(err)

			return
		}

		message := string(b)
		if message == "Hello" {
			if err := connection.Send([]byte(message + ", World!")); err != nil {
				log.Print(err)

				return
			}
		}
	})
	if err != nil {
		log.Println(err)
	}

	log.Fatal(server.Run())
}
