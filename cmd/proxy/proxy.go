package main

import (
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

func get_server_conn() (net.Conn, error) {
	s_addr, err := net.ResolveTCPAddr(conn_type, server_addr)
	if err != nil {
		return nil, fmt.Errorf("Could not resolve server address: %s\n", err.Error())
	}

	conn, err := net.DialTCP(conn_type, nil, s_addr)
	if err != nil {
		return nil, fmt.Errorf("Connection refused by server: %s\n", err.Error())
	}

	return conn, nil
}

func handle_conn(client_conn net.Conn) {
	defer client_conn.Close()

	addr := client_conn.RemoteAddr().String()
	defer log.Printf("Closing socket %s\n", addr)

	server_conn, err := get_server_conn()
	defer server_conn.Close()

	if err != nil {
		log.Fatalf("Error connecting to server: %s\n", err.Error())
	}

	log.Printf("Forwarding request to server at address %s\n", server_addr)

	_, err = io.Copy(server_conn, client_conn)
	if err != nil {
		log.Fatalf("Error forwarding request to server: %s\n", err.Error())
	}

	_, err = io.Copy(client_conn, server_conn)
	if err != nil {
		log.Fatalf("Error sending response to client: %s\n", err.Error())
	}
}

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
