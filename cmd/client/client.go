package main

import (
	"io"
	"log"
	"net"
)

func main() {
	server_addr, err := net.ResolveTCPAddr("tcp4", ":7324")
	if err != nil {
		log.Fatalf("Error issued when resolving server address: %s\n", err.Error())
	}

	conn, err := net.DialTCP("tcp4", nil, server_addr)
	if err != nil {
		log.Fatalf("Error issued when connecting to server: %s\n", err.Error())
	}
	defer conn.Close()

	n, err := conn.Write([]byte("Hello, World!\r\n"))
	if err != nil {
		log.Fatalf("Error issued when writing buffer: %s\n", err.Error())
	}

	log.Printf("Total bytes sent: %d", n)

	buf := make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Fatalf("Error issued when reading response: %s\n", err.Error())
		}
	}

	log.Printf("Got %d bytes\n", len(buf))
	log.Printf("Response received from the server: %s\n", string(buf))
}
