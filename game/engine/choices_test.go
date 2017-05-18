package engine

import (
	"testing"
)

func TestChoicesWinTable(t *testing.T) {
	var choices []Choice = []Choice{Rock, Paper, Scissors, Lizard, Spock}
	for i := range choices {
		if choices[i].Beats(choices[i]) {
			t.Errorf("One of the choices (%s) beats itself", choices[i].Name())
		}
	}

	if Rock.Beats(Paper) || Rock.Beats(Spock) ||
		Paper.Beats(Scissors) || Paper.Beats(Lizard) ||
		Scissors.Beats(Rock) || Scissors.Beats(Spock) ||
		Lizard.Beats(Rock) || Lizard.Beats(Scissors) ||
		Spock.Beats(Paper) || Spock.Beats(Paper) {
		t.Errorf("At least one of the choices beats something it should not.")
	}

	if !Rock.Beats(Scissors) || !Rock.Beats(Lizard) ||
		!Paper.Beats(Rock) || !Paper.Beats(Spock) ||
		!Scissors.Beats(Paper) || !Scissors.Beats(Lizard) ||
		!Lizard.Beats(Paper) || !Lizard.Beats(Spock) ||
		!Spock.Beats(Rock) || !Spock.Beats(Scissors) {
		t.Errorf("At least one of the choices does not beat something it should.")
	}
}
