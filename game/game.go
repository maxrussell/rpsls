package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/maxrussell/rpsls"
	"github.com/maxrussell/rpsls/game/engine"
)

func main() {
	http.HandleFunc("/choices", getChoices)
	http.HandleFunc("/choice", getRandomChoice)
	http.HandleFunc("/play", play)
	http.HandleFunc("/health", getHealthCheck)

	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

func getChoices(response http.ResponseWriter, _ *http.Request) {
	choices := engine.GetChoices()
	choicesJson, err := json.Marshal(choices)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	response.Write(choicesJson)
}

func getRandomChoice(response http.ResponseWriter, _ *http.Request) {
	choice, err := engine.GetRandomChoice()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	choiceJson, err := json.Marshal(choice)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	response.Write(choiceJson)
}

func play(response http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()

	player, playerSet := queryParams["player"]
	playerChoiceQuery, choiceSet := queryParams["choice"]
	computer, computerSet := queryParams["computer"]

	if !playerSet || len(player) < 1 || len(player[0]) < 1 || !computerSet || len(computer) < 1 || len(computer[0]) < 1 {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("The player and computer query parameters (and values) are required and may not be empty."))
		return
	}

	computerChoice, err := engine.GetRandomChoice()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Unable to choose randomly for the computer: " + err.Error()))
		return
	}

	var playerChoice engine.Choice
	if !choiceSet || len(playerChoiceQuery) < 1 {
		playerChoice, err = engine.GetRandomChoice()
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("Unable to choose randomly for the player: " + err.Error()))
			return
		}
	} else {
		parsedPlayerChoice, err := strconv.ParseInt(playerChoiceQuery[0], 10, 32)
		if err != nil {
			response.WriteHeader(422)
			response.Write([]byte("Player choice was malformed."))
			return
		}

		playerChoice = engine.Choice(parsedPlayerChoice)
		if !playerChoice.Valid() {
			response.WriteHeader(422)
			response.Write([]byte("Invalid choice made"))
			return
		}
	}

	result := struct {
		Results string
	}{}
	if playerChoice == computerChoice {
		result.Results = "tie"
	} else if playerChoice.Beats(computerChoice) {
		result.Results = "win"
	} else {
		result.Results = "lose"
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	response.Write(resultJson)
}

func getHealthCheck(response http.ResponseWriter, _ *http.Request) {
	healthJson, err := json.Marshal(rpsls.HealthCheck{
		Status:    "OK",
		Name:      rpsls.HostName(),
		Version:   "1.0",
		StartTime: rpsls.StartTimeString(),
		UpTime:    rpsls.UpTimeString(),
	})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	response.Write(healthJson)
}
