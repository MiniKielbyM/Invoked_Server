package main

var cards = map[string]card{
	"CITIZEN": {
		Name:           "Citizen",
		Description:    "The common rabble.",
		Suit:           "ROSE",
		BufferTurnsMax: 1,
		Cost:           cost{
			Chips: 1,
			Burn:  0,
			Sack:  0,
		},
		BaseStats:      statBlock{
			Attack:  75,
			Defense: 100,
		},
		Type:  Conscript,
	},
	"WIRE_ABOMINATON": {
		Name:           "Wire Abomination",
		Description:    "One of few half-stable successes.",
		Suit:           "IVORY",
		BufferTurnsMax: 1,
		Cost:           cost{
			Chips: 6,
			Burn:  0,
			Sack:  0,
		},
		BaseStats:      statBlock{
			Attack:  400,
			Defense: 200,
		},
		Type:  Conscript,
	},
	"STRENGTHENED": {
		Name:           "Strengthened",
		Description:    "Somtimes the unknowns gifts are simple.",
		Suit:           "PROLIFERATED",
		BufferTurnsMax: 3,
		Cost:           cost{
			Chips: 2,
		},
		BaseStats:      statBlock{
			Attack:  200,
			Defense: 150,
		},
		Type:  Conscript,
	},
	"SOLDIER": {
		Name:           "Soldier",
		Description:    "Old and brave, or young and stupid. It's a coin toss really.",
		Suit:           "IRON",
		BufferTurnsMax: 1,
		Cost:           cost{
			Chips: 1,
		},
		BaseStats:      statBlock{
			Attack:  100,
			Defense: 75,
		},
		Type:  Conscript,
	},
	"DRUNKARD": {
		Name:           "Drunkard",
		Description:    "With how common they are you'd think liqour was legal.",
		Suit:           "LOTTERY",
		BufferTurnsMax: 1,
		Cost:           cost{
			Chips: 1,
		},
		BaseStats:      statBlock{
			Attack:  25,
			Defense: 125,
		},
		Type:  Conscript,
	},

}

