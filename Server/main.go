package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)



func main() {
	// Start listening on TCP
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Listen failed:", err)
	}
	fmt.Println("Server is listening on localhost:8080" )
	defer listener.Close()

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		// Handle each connection in a goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr()

	reader := bufio.NewReader(conn)
	for {
		// Read client message
		raw, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %v", err)
			}
			return
		}

		raw = strings.TrimSuffix(raw, "\n")

		// Parse headers and message (format: "header1,header2,...,message")
		parts := strings.Split(raw, "||HEADER.END||")

		var headers []string
		var message string
		if len(parts) > 1 {
			headers = strings.Split(strings.Split(parts[0], "||HEADER.START||")[1], ",")
			message = parts[1]
		} else {
			message = raw
		}

		fmt.Printf("Got message from %s\n", clientAddr)
		fmt.Printf("  Headers: %v\n", headers)
		fmt.Printf("  Message: %s\n", message)

		// Echo back
		_, err = conn.Write([]byte(raw + "\n"))
		if err != nil {
			log.Printf("Write error: %v", err)
			return
		}
	}
}
