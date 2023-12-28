package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/oleh-malakan/go-node"
)

func Handler(listener *node.Listener, nodeID string, f func(stream *node.Stream)) {
	go func() {
		for {
			stream, err := listener.Accept()
			if err != nil {
				break
			}

			go func() {
				f(stream)
				stream.Close()
			}()
		}
	}()
}

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	server, err := node.Run(0, &net.UDPAddr{})
	if err != nil {
		log.Fatal(err)
	}

	listener, err := server.Listen("Hello, World!")
	if err != nil {
		log.Print(err)

		return
	}

	Handler(listener, "Hello, World!", func(stream *node.Stream) {
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

	<-exit
}
