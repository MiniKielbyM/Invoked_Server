package main

import "net"

// Server structs

type message struct {
	Headers []string
	Body    string
}

// Game substructs

type cost struct {
	Chips int
	Burn int
	Sack int
}

type statBlock struct {
	Attack    int
	Defense   int
}

type witness struct {
	Turn 	int
	Time   	int
	Victim 	*card
	Killer 	*card
	Card   	*card
}

// Game structs

type card struct {
	Name        	string
	Description 	string
	Suit 	    	string
	BufferTurns 	int
	BufferTurnsMax 	int
	Snuffed   		bool
	Location  		cardLocation
	Cost        	cost
	BaseStats    	statBlock // DONT TOUCH THESE
	MutableStats 	statBlock
	Type         	cardType
	Owner   		*player
	Controller  	*player
	Witnesses   	[]witness
}

func (c *card) ResetStats() {

}

// Game superstructs

type player struct {
	Conn  *net.Conn
	Hand  hand
	Deck  deck
	Id    string
	Name  string
	Chips int
}
type deck struct {
	Cards []*card
}
type hand struct {
	Cards []*card
}
type table struct {
	Player1 player
	Player2 player
}

// Game main struct

type game struct {
	Id 		string
	Table   table
}