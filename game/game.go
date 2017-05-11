package main

import (
	"encoding/json"
	"log"
	"net/http"

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

func play(response http.ResponseWriter, _ *http.Request) {

}

func getHealthCheck(response http.ResponseWriter, _ *http.Request) {

}
