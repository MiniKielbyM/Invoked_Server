package main

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
)

// CryptoRandSource implements rand.Source to use crypto/rand for secure randomness.
type CryptoRandSource struct{}

func (s CryptoRandSource) Int63() int64 {
	return int64(s.Uint64() & (1<<63 - 1)) // mask off sign bit
}

func (s CryptoRandSource) Uint64() uint64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Fatal("crypto/rand error:", err)
	}
	return binary.LittleEndian.Uint64(b[:])
}

func (s CryptoRandSource) Seed(seed int64) {} // Seed is a no-op for crypto/rand

// Reserved headers that only the server can send
var reservedHeaders = map[string]bool{
	"LobbyCreate":      true,
	"LobbyJoin":        true,
	"LobbyList":        true,
	"GameStart":        true,
	"PlayerConnect":    true,
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

func (d *deck) shuffle() {
	// Simple shuffle implementation using math/rand with crypto/rand source
	source := CryptoRandSource{}
	for i := len(d.Cards) - 1; i > 0; i-- {
		j := int(source.Uint64() % uint64(i+1))
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
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
	if len(message.Body) == 0 {
		return fmt.Errorf("lobby name cannot be empty")
	}
	if players[message.Sender].CurrentLobbyId != "" {
		return fmt.Errorf("player is already in a lobby")
	}
	lobbyName := message.Body
	lobbyId := fmt.Sprintf("lobby%d", len(lobbies)+1)

	newLobby := &lobby{
		Id:   lobbyId,
		Name: lobbyName,
		Host: players[message.Sender],
		Players: []player{
			players[message.Sender],
		},
		InProgress: false,
	}
	lobbies[lobbyId] = newLobby
	p := players[message.Sender]
	p.CurrentLobbyId = newLobby.Id
	players[message.Sender] = p
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
	if players[message.Sender].CurrentLobbyId != "" {
		return fmt.Errorf("player is already in a lobby")
	}
	p := players[message.Sender]
	p.CurrentLobbyId = lobby.Id
	players[message.Sender] = p
	lobby.Players = append(lobby.Players, players[message.Sender])
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
		Conn: &message.Sender,
		Deck: deck{
			Cards: parseDeckString(strings.Split(strings.Split(message.Body, "||PLAYER.DECK||")[1], "||")[0]),
		},
		Id:   string(newUUID),
		Name: strings.Split(strings.Split(message.Body, "||PLAYER.NAME||")[1], "||")[0],
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
