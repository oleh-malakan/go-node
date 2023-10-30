package main

import (
	"fmt"
	"log"

	"github.com/oleh-malakan/go-node"
)

func main() {
	client, err := node.Dial()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := client.Connect("Hello, World!", nil)
	if err != nil {
		log.Print(err)

		return
	}

	if err := conn.Send([]byte("Hello")); err != nil {
		log.Print(err)

		return
	}

	b, err := conn.Receive()
	if err != nil {
		log.Print(err)

		return
	}
	fmt.Println(string(b))

	if err := conn.Close(); err != nil {
		log.Print(err)

		return
	}
}
