package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// Card Pool - Available cards for deck building


// ValidateCard checks if a card exists in the pool
func ValidateCard(cardName string) (card, bool) {
	card, exists := cards[cardName]
	return card, exists
}

// ValidateDeck checks if all cards in a deck list exist in the pool
func ValidateDeck(cardNames []string) (bool, []string) {
	var invalidCards []string
	for _, name := range cardNames {
		if _, exists := ValidateCard(name); !exists {
			invalidCards = append(invalidCards, name)
		}
	}
	return len(invalidCards) == 0, invalidCards
}

func parseMessage(raw string) message {
	parts := strings.Split(raw, "||HEADER.END||")

	var headers []string
	if len(parts) > 1 {
		headers = strings.Split(parts[0], "||HEADER.SEP||")
	}

	body := ""
	if len(parts) > 1 {
		body = parts[1]
	} else {
		body = parts[0]
	}

	return message{
		Headers: headers,
		Body:    body,
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
		msg := parseMessage(raw)
		headers := msg.Headers
		message := msg.Body
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
