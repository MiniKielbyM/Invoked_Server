package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	// Set up UDP address
	addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		log.Fatal("Couldn't resolve address:", err)
	}

	// Start listening
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("Listen failed:", err)
	} else {
		fmt.Println("Server is listening on localhost:8080")
	}
	defer conn.Close()

	// Buffer for incoming data
	buffer := make([]byte, 1024)
	for {
		// Read client message
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Read error: %v", err)
			continue
		}

		// Parse headers and message (format: "header1,header2,...,message")
		raw := string(buffer[:n])
		parts := strings.Split(raw, "||HEADER.END||")

		var headers []string
		var message string
		if len(parts) > 1 {
			headers = parts[:0]
			message = parts[1]
		} else {
			message = raw
		}

		fmt.Printf("Got message from %s\n", clientAddr)
		fmt.Printf("  Headers: %v\n", headers)
		fmt.Printf("  Message: %s\n", message)

		// Echo back
		_, err = conn.WriteToUDP(buffer[:n], clientAddr)
		if err != nil {
			log.Printf("Write error: %v", err)
		}
	}
}
