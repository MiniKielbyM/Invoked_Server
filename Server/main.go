package main

import (
	"errors"
	"fmt"
	"log"
	"net"
)


func main() {
	// Start listening on TCP
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Listen failed:", err)
	}
	fmt.Println("Server is listening on localhost:8080")
	defer listener.Close()

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				fmt.Println("Listener closed; server shutting down.")
				break
			}
			log.Printf("Accept error: %v", err)
			continue
		}

		// Handle each connection in a goroutine
		go handleConnection(conn)
	}
}


