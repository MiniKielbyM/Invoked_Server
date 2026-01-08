package main


// Game enums

// Card Types
type cardType string
const (
	Conscript cardType = "Conscript"
	Joker     cardType = "Joker"
	Strategy  cardType = "Strategy"
	Building  cardType = "Building"
	Action    cardType = "Action"
)

type cardLocation string
const (
	Hand      	cardLocation = "Hand"
	Deck      	cardLocation = "Deck"
	Table     	cardLocation = "Table"
	Graveyard 	cardLocation = "Graveyard"
	Ashtray   	cardLocation = "Ashtray"
	Seal      	cardLocation = "Seal"
)
