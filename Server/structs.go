package main

import "net"

// Server structs

type message struct {
	Sender  net.Conn
	Headers []string
	Body    string
}

type lobby struct {
	Id         string
	Name       string
	InProgress bool
	Host       player
	Players    []player
	Game       game
}

// Game substructs

type cost struct {
	Chips    int
	Burn     int
	Sack     int
	Vitality int
}

type statBlock struct {
	Attack  int
	Defense int
}

type witness struct {
	Turn   int
	Time   int
	Victim *card
	Killer *card
	Card   *card
}

// Game structs

type card struct {
	Name           string
	Description    string
	Suit           string
	BufferTurns    int
	BufferTurnsMax int
	Snuffed        bool
	Location       cardLocation
	Cost           cost
	BaseStats      statBlock // DONT TOUCH THESE
	MutableStats   statBlock
	Type           cardType
	Owner          *player
	Controller     *player
	Witnesses      []witness
}

func (c *card) ResetStats() {

}

// Game superstructs

type player struct {
	Conn    *net.Conn
	Hand    hand
	Deck    deck
	Discard deck
	Ashtray deck
	Id      string
	Name    string
	Chips   int
	Health  int
	CurrentLobbyId string
}
type deck struct {
	Cards []card
}
type hand struct {
	Cards []card
}
type table struct {
	Player1    player
	P1FrontRow [][]card
	P1BackRow  [][]card
	Player2    player
	P2FrontRow [][]card
	P2BackRow  [][]card
	SealZone   [][]card
}

// Game main struct

type game struct {
	Id    string
	Table table
}
