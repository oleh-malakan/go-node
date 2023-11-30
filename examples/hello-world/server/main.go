package main

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/oleh-malakan/go-node"
)

func main() {
	server, err := node.New(&tls.Config{})
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
		log.Print(err)

		return
	}

	log.Fatal(server.Run(&net.UDPAddr{}, 0))
}
