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

	_, err = server.Handler("Hello, World!", func(stream *node.NodeStream) {
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
	})
	if err != nil {
		log.Print(err)

		return
	}

	log.Fatal(server.Run(0))
}
