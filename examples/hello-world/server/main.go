package main

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/oleh-malakan/go-node"
)

func main() {
	node.Handler("Hello, World!", func(query []byte, connection *node.Connection) {
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

	log.Fatal(node.Do(&tls.Config{}, 9999, &net.UDPAddr{}))
}
