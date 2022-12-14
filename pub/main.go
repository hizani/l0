package main

import (
	"fmt"
	"os"

	"github.com/nats-io/stan.go"
)

func main() {
	var text []byte
	var conn stan.Conn
	var err error

	// Get filename as an argument
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%v", "pub [FILE]")
		return
	}
	filename := &(os.Args[1])

	// Establish connection with nats-streaming server
	conn, err = stan.Connect("test-cluster", "publisher", stan.NatsURL("127.0.0.1:4222"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}
	defer conn.Close()

	// Read file
	text, err = os.ReadFile(*filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	// Publish to the channel
	if err = conn.Publish("chan", text); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}
}
