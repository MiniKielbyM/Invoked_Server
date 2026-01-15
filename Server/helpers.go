package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"os/exec"
)

// Reserved headers that only the server can send
var reservedHeaders = map[string]bool{
	"LobbyCreate": true,
	"LobbyJoin":   true,
	"LobbyList":   true,
	"GameStart":   true,
	"PlayerConnect": true,
	"PlayerUpdateDeck": true,
}

// isReservedHeader checks if any of the message headers are reserved for server use
func isReservedHeader(headers []string) bool {
	for _, header := range headers {
		if reservedHeaders[header] {
			return true
		}
	}
	return false
}

// ValidateCard checks if a card exists in the pool
func ValidateCard(cardName string) (card, bool) {
	card, exists := cards[cardName]
	return card, exists
}

func parseDeckString(deckStr string) []card {
	if deckStr == "" {
		return []card{}
	}
	cardNames := strings.Split(deckStr, ",")
	var deck []card
	for _, name := range cardNames {
		if card, exists := ValidateCard(name); exists {
			deck = append(deck, card)
		}
	}
	return deck
}

func parseMessage(raw string) message {
	// Strip UTF-8 BOM if present
	raw = strings.TrimPrefix(raw, "\xEF\xBB\xBF")

	parts := strings.Split(raw, "||HEADER.END||")

	var headers []string
	if len(parts) > 1 {
		headers = strings.Split(parts[0], "||HEADER.SEP||")
		// Strip BOM from each header if present
		for i := range headers {
			headers[i] = strings.TrimPrefix(headers[i], "\xEF\xBB\xBF")
		}
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

func sendMessage(conn net.Conn, headers []string, body string) error {
	message := strings.Join(headers, "||HEADER.SEP||") + "||HEADER.END||" + body + "\n"
	_, err := conn.Write([]byte(message))
	return err
}

// Handlers
func handleLobbyCreate(message message) error {
	lobbyName := message.Body
	lobbyId := fmt.Sprintf("lobby%d", len(lobbies)+1)
	newLobby := &lobby{
		Id:      lobbyId,
		Name:    lobbyName,
		Host:    player{
			Conn:  &message.Sender,
			Hand:  hand{},
			Deck:  deck{},
			Id:    "",
			Name:  "",
			Chips: 0,
		},
	}
	lobbies[lobbyId] = newLobby
	fmt.Printf("Created lobby: %s with ID: %s\n", lobbyName, lobbyId)
	sendMessage(message.Sender, []string{"LobbyCreate"}, fmt.Sprintf("Lobby %s created with ID %s", lobbyName, lobbyId))
	return nil
}

func handleLobbyJoin(message message) error {
	lobbyId := message.Body
	lobby, exists := lobbies[lobbyId]
	if !exists {
		return fmt.Errorf("lobby with ID %s does not exist", lobbyId)
	}
	if lobby.InProgress {
		return fmt.Errorf("cannot join lobby %s: game already in progress", lobbyId)
	}
	if len(lobby.Players) == 2 {
		return fmt.Errorf("cannot join lobby %s: lobby is full", lobbyId)
	}
	lobby.Players = append(lobby.Players, player{
		Conn:  &message.Sender,
		Hand:  hand{},
		Deck:  deck{},
		Id:    "",
		Name:  "",
		Chips: 0,
	})
	fmt.Printf("Player joined lobby: %s\n", lobbyId)
	sendMessage(message.Sender, []string{"LobbyJoin"}, fmt.Sprintf("Joined lobby %v", lobby))
	return nil
}

func handleLobbyList(message message) error {
	fmt.Println("Listing lobbies:")
	for id, lobby := range lobbies {
		fmt.Printf("Lobby ID: %s, Name: %s, Players: %d\n, In progress: %t\n", id, lobby.Name, len(lobby.Players), lobby.InProgress)
	}
	sendMessage(message.Sender, []string{"LobbyList"}, fmt.Sprintf("%v", lobbies))
	return nil
}

func handlePlayerConnect(message message) error {
	newUUID, _ := exec.Command("uuidgen").Output()
	players[message.Sender] = player{
		Conn:  &message.Sender,
		Deck:  deck{
			Cards: parseDeckString(strings.Split(strings.Split(message.Body, "||PLAYER.DECK||")[1], "||")[0]),
		},
		Id:    string(newUUID),
		Name:  strings.Split(strings.Split(message.Body, "||PLAYER.NAME||")[1], "||")[0],
	}

	fmt.Printf("Player connected: %v\n", players[message.Sender])
	sendMessage(message.Sender, []string{"PlayerConnect"}, fmt.Sprintf("Player %s connected", string(newUUID)))
	return nil
}

func handlePlayerUpdateDeck(message message) error {
	player, exists := players[message.Sender]
	if !exists {
		return fmt.Errorf("player not found for connection")
	}
	player.Deck = deck{
		Cards: parseDeckString(message.Body),
	}
	players[message.Sender] = player
	fmt.Printf("Player %s updated deck\n", player.Id)
	sendMessage(message.Sender, []string{"PlayerUpdateDeck"}, "Deck updated successfully")
	return nil
}

func handleGameStart(message message) error {
	for lobbyId, lobby := range lobbies {
		if message.Sender == *lobby.Host.Conn {
			if len(lobby.Players) < 2 {
				return fmt.Errorf("not enough players to start the game in lobby %s", lobbyId)
			}
			if lobby.InProgress {
				return fmt.Errorf("game in lobby %s is already in progress", lobbyId)
			}
			fmt.Printf("Starting game in lobby: %s\n", lobbyId)
			lobby.InProgress = true
			sendMessage(*lobby.Players[0].Conn, []string{"GameStart"}, "Game has started")
			sendMessage(*lobby.Players[1].Conn, []string{"GameStart"}, "Game has started")
			player := &lobby.Players[0]
			player2 := &lobby.Players[1]
			lobby.Game = game{
				Table: table{
					Player1: *player,
					Player2: *player2,
				},
			}
			return nil
		}
	}
	return fmt.Errorf("sender is not the host of any lobby")
}

// Updated handleConnection
func handleConnection(conn net.Conn, handler *Handler) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr()

	reader := bufio.NewReader(conn)
	for {
		raw, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %v", err)
			}
			return
		}

		raw = strings.TrimSuffix(raw, "\n")
		msg := parseMessage(raw)

		fmt.Printf("Got message from %s\n", clientAddr)
		fmt.Printf("  Headers: %v\n", msg.Headers)
		fmt.Printf("  Message: %s\n", msg.Body)

		// Check if client is trying to use reserved headers
		if isReservedHeader(msg.Headers) {
			sendMessage(conn, []string{"Error"}, "Cannot use reserved headers")
			continue
		}

		// Route the message
		if err := handler.Route(message{
			Sender:  conn,
			Headers: msg.Headers,
			Body:    msg.Body,
		}); err != nil {
			log.Printf("Handler error: %v", err)
		}
	}
}
