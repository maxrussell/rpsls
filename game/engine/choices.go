package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	Rock     = 1
	Paper    = 2
	Scissors = 3
	Lizard   = 4
	Spock    = 5
)

var names map[Choice]string = map[Choice]string{
	Rock:     "Rock",
	Paper:    "Paper",
	Scissors: "Scissors",
	Lizard:   "Lizard",
	Spock:    "Spock",
}

type Choice int

func (c Choice) Name() string {
	return names[c]
}

func (c Choice) Beats(choice Choice) bool {
	for _, thingCBeats := range winTable[c] {
		if thingCBeats == choice {
			return true
		}
	}
	return false
}

func (c Choice) Valid() bool {
	for _, choice := range GetChoices() {
		if c == choice {
			return true
		}
	}
	return false
}

type choiceJson struct {
	Id   int
	Name string
}

func (c Choice) MarshalJSON() ([]byte, error) {
	return json.Marshal(choiceJson{
		Id:   int(c),
		Name: c.Name(),
	})
}

func GetChoices() []Choice {
	return []Choice{1, 2, 3, 4, 5}
}

func GetRandomChoice() (Choice, error) {
	choices := GetChoices()
	randomNumber, err := getRandomNumber()
	if err != nil {
		return 0, fmt.Errorf("error generating a random number: %s", err.Error())
	}

	return choices[randomNumber%len(choices)], nil
}

// Consider abstracting this for testability
func getRandomNumber() (int, error) {
	response, err := http.Get("http://codechallenge.boohma.com/random")
	if err != nil {
		return -1, err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return -1, err
	}

	decodedResponse := struct {
		RandomNumber int `json:"random_number"`
	}{}
	err = json.Unmarshal([]byte(responseBody), &decodedResponse)
	if err != nil {
		return -1, err
	}

	return decodedResponse.RandomNumber, nil
}
