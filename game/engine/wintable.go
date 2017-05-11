package engine

// Map from a choice to what it beats
var winTable = map[Choice][]Choice{
	Rock:     []Choice{Scissors, Lizard},
	Paper:    []Choice{Rock, Spock},
	Scissors: []Choice{Paper, Lizard},
	Lizard:   []Choice{Paper, Spock},
	Spock:    []Choice{Rock, Scissors},
}
