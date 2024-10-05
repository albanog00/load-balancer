package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

const (
	conn_type  = "tcp4"
	proxy_addr = "127.0.0.1:7324"
)

var server_addrs = []string{
	"127.0.0.1:7325",
	"127.0.0.1:7326",
	"127.0.0.1:7327",
	"127.0.0.1:7328",
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

	conn.Write(builder.Bytes())
}

func main() {
	wait := sync.WaitGroup{}
	wait.Add(len(server_addrs))

	signal_chan := make(chan os.Signal, len(server_addrs))
	signal.Notify(signal_chan, os.Interrupt)

	go func() {
		<-signal_chan
		log.Printf("Interrupt signal received. Closing listeners.\n")
		for range server_addrs {
			wait.Done()
		}
	}()

	for _, server_addr := range server_addrs {
		go func() {
			listener, err := net.Listen(conn_type, server_addr)
			if err != nil {
				log.Fatalf("Error issued: %s\n", err.Error())
			}

			log.Printf("Listening on %s\n", server_addr)

			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Fatalf("Error while accepting connection: %s\n", err.Error())
					break
				}

				go handle_conn(conn)
			}
		}()
	}

	wait.Wait()
}
