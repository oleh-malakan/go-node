package main

import (
	"log"

	"github.com/oleh-malakan/go-node"
)

func main() {
	node.Handler("Hello, World!", func(c *node.Connection) {
		b, err := c.Receive()
		if err != nil {
			log.Print(err)

			return
		}

		msg := string(b)
		if msg == "Hello" {
			if err := c.Send([]byte(msg + ", World!")); err != nil {
				log.Print(err)

				return
			}
		}
	})

	log.Fatal(node.ListenAndServe())
}
