package main

import (
	"io"
	"log"
	"net"
	"sync"
)

func main() {
	server_addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:7324")
	if err != nil {
		log.Fatalf("Error issued when resolving server address: %s\n", err.Error())
	}

	wait := sync.WaitGroup{}
	successful := 0

	for range 100 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for range 1000 {
				conn, err := net.DialTCP("tcp4", nil, server_addr)
				if err != nil {
					log.Fatalf("Error issued when connecting to server: %s\n", err.Error())
				}

				defer conn.Close()

				n, err := conn.Write([]byte("Hello, World!\r\n"))
				if err != nil {
					log.Fatalf("Error issued when writing buffer: %s\n", err.Error())
				}

				// Connection blocks on read if write don't get closed ??
				conn.CloseWrite()
				log.Printf("Total bytes sent: %d", n)

				buf := make([]byte, 1024)
				n, err = conn.Read(buf)
				if err != nil {
					if err != io.EOF {
						log.Fatalf("Error issued when reading response: %s\n", err.Error())
					}
				}

				log.Printf("Got %d bytes\n", n)
				log.Printf("Response received from the server: %s\n", string(buf[:n]))

				successful++
			}
		}()
	}

	wait.Wait()
	log.Printf("\nCount of successful connections: %d\n", successful)
}
