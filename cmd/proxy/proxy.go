package main

import (
	"bytes"
	"fmt"
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
	listener, err := net.Listen(conn_type, proxy_addr)
	if err != nil {
		log.Fatalf("Error issued: %s\n", err.Error())
	}

	log.Printf("Listening on %s\n", proxy_addr)

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

	// forward request to origin server
	log.Printf("Forwarding request to server at address %s\n", server_addr)
	response, err := send_req_to_server(builder.Bytes())
	if err != nil {
		log.Fatalf("Error issued when forwarding request to server: %s\n", err.Error())
	}

	conn.Write(response)
}

func send_req_to_server(content []byte) ([]byte, error) {
	s_addr, err := net.ResolveTCPAddr(conn_type, server_addr)
	if err != nil {
		return []byte{}, fmt.Errorf("Could not resolve server address: %s\n", err.Error())
	}

	conn, err := net.DialTCP(conn_type, nil, s_addr)
	if err != nil {
		return []byte{}, fmt.Errorf("Connection refused by server: %s\n", err.Error())
	}
	defer conn.Close()

	builder := bytes.Buffer{}
	buf := make([]byte, 1024)

	n, err := conn.Write(content)
	if err != nil {
		return []byte{}, fmt.Errorf("Unable to write response from the server: %s\n", err.Error())

	}

	n, err = conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			return []byte{}, fmt.Errorf("Unable to read response from the server: %s\n", err.Error())
		}
	}

	builder.Write(buf[:n])

	return builder.Bytes(), nil
}
