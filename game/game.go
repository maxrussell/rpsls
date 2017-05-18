package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/maxrussell/rpsls"
	"github.com/maxrussell/rpsls/game/engine"
)

var scoreboardLocation = os.Getenv("SCOREBOARD")

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/choices", getChoices)
	mux.HandleFunc("/choice", getRandomChoice)
	mux.HandleFunc("/play", play)
	mux.HandleFunc("/health", getHealthCheck)

	err := http.ListenAndServe(":8080", rpsls.DefaultToPlainText(mux))
	log.Fatal(err)
}

func getChoices(response http.ResponseWriter, _ *http.Request) {
	choices := engine.GetChoices()
	choicesJson, err := json.Marshal(choices)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Error marshaling available choices as JSON."))
		log.Println("Error marshaling available choices as JSON:", err.Error())
		return
	}

	rpsls.WriteAsJson(response, choicesJson)
}

func getRandomChoice(response http.ResponseWriter, _ *http.Request) {
	choice, err := engine.GetRandomChoice()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Error getting random choice."))
		log.Println("Error getting random choice:", err.Error())
		return
	}

	choiceJson, err := json.Marshal(choice)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Error marshaling random choice as JSON."))
		log.Println("Error marshaling random choice as JSON:", err.Error())
		return
	}

	rpsls.WriteAsJson(response, choiceJson)
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
		response.Write([]byte("Unable to choose randomly for the computer."))
		log.Println("Error trying to choose randomly for the computer:", err.Error())
		return
	}

	var playerChoice engine.Choice
	if !choiceSet || len(playerChoiceQuery) < 1 {
		playerChoice, err = engine.GetRandomChoice()
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("Unable to choose randomly for the player."))
			log.Println("Error trying to choose randomly for the player:", err.Error())
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

	if len(scoreboardLocation) > 0 {
		var playerScore, computerScore int
		if result.Results == "win" {
			playerScore = 1
		} else if result.Results == "lose" {
			computerScore = 1
		}

		players := []rpsls.Player{
			rpsls.Player{
				UserName: player[0],
				Score:    playerScore,
			},
			rpsls.Player{
				UserName: computer[0],
				Score:    computerScore,
			},
		}
		playersJson, err := json.Marshal(players)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("Error trying to marshal results as JSON; results were not posted to scoreboard."))
			log.Println("Error trying to marshal results as JSON; not posting to scoreboard. Error:", err.Error())
			return
		}
		body := strings.NewReader(string(playersJson))
		_, err = http.Post(scoreboardLocation+"/results", "application/json", body)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("There was an error sending results to the scoreboard."))
			log.Println("Error sending results to scoreboard:", err.Error())
			return
		}
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Error marshaling results as JSON."))
		log.Println("Error marshaling results as JSON:", err.Error())
		return
	}

	rpsls.WriteAsJson(response, resultJson)
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
		response.Write([]byte("Error marshaling as JSON."))
		log.Println("Error marshaling health check results as JSON:", err.Error())
		return
	}

	rpsls.WriteAsJson(response, healthJson)
}
