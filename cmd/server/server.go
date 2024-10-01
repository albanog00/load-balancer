package main

import (
	"bytes"
	"io"
	"log"
	"net"
)

const (
	conn_type   = "tcp4"
	server_addr = "127.0.0.1:7325"
	proxy_addr  = "127.0.0.1:7324"
)

func main() {
	listener, err := net.Listen(conn_type, server_addr)
	if err != nil {
		log.Fatalf("Error issued: %s\n", err.Error())
	}

	log.Printf("Listening on %s\n", server_addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error issued while accepting connection: %s\n", err.Error())
		}

		go handle_conn(conn)
	}
}

func handle_conn(conn net.Conn) {
	addr := conn.RemoteAddr().String()

	defer conn.Close()
	defer log.Printf("Closing socket %s\n", addr)

	// if addr != proxy_addr {
	// 	log.Printf("%s unauthorized\n", addr)
	// 	return
	// }

	builder := bytes.Buffer{}
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Fatalf("Error issued while reading from connection %s: %s\n", addr, err.Error())
		}
	}

	builder.Write(buf[:n])

	log.Printf("Got %d bytes\n", len(builder.Bytes()))
	log.Printf("Request from %s\n%s\n", addr, builder.String())

	conn.Write(builder.Bytes())
}
