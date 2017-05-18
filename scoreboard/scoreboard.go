package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/maxrussell/rpsls"
	"github.com/maxrussell/rpsls/scoreboard/storage"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/scoreboard", getScoreboard)
	mux.HandleFunc("/results", addResult)
	mux.HandleFunc("/health", getHealthCheck)

	err := http.ListenAndServe(":8081", rpsls.DefaultToPlainText(mux))
	log.Fatal(err)
}

func getScoreboard(response http.ResponseWriter, request *http.Request) {
	topPlayersToShow := 10
	query := request.URL.Query()
	if countStr := query.Get("count"); len(countStr) > 0 {
		count, err := strconv.ParseInt(countStr, 10, 32)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("Parameter count must be a base-ten integer or omitted."))
			return
		}

		topPlayersToShow = int(count)
	}

	players, err := storage.GetTopPlayers(topPlayersToShow)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Error retrieving top players."))
		log.Println("Error getting top players from storage:", err.Error())
		return
	}

	playersJson, err := json.Marshal(players)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Error marshaling response as JSON."))
		log.Println("Error marshaling top players as JSON:", err.Error())
		return
	}

	rpsls.WriteAsJson(response, playersJson)
}

func addResult(response http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	defer func() {
		if nil == body {
			return
		}
		request.Body.Close()
	}()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Error reading request body."))
		log.Println("Error reading request body:", err.Error())
		return
	}

	players := []rpsls.Player{}
	err = json.Unmarshal(body, &players)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Request body was malformed and couldn't be unmarshaled as JSON."))
		return
	}

	if len(players) > 2 || len(players) < 1 {
		response.WriteHeader(422) // not valid
		response.Write([]byte("Only two players may participate in a game."))
		return
	}

	if players[0].Score > 0 && players[1].Score > 0 {
		response.WriteHeader(422)
		response.Write([]byte("Only one player may win a game."))
		return
	}

	if players[0].UserName == players[1].UserName {
		response.WriteHeader(422)
		response.Write([]byte("A player may not play his- or herself."))
		return
	}

	var winner, loser rpsls.Player
	if players[0].Score > 0 {
		winner = players[0]
		loser = players[1]
	} else {
		loser = players[0]
		winner = players[1]
	}

	err = storage.AddResult(winner, loser)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Error adding result to storage."))
		log.Println("Error adding result to storage:", err.Error())
		return
	}

	response.Header().Del("Content-Type")
	response.WriteHeader(http.StatusNoContent)
}

func getHealthCheck(response http.ResponseWriter, request *http.Request) {
	healthCheck := storage.GetHealthCheck()
	healthCheckJson, err := json.Marshal(healthCheck)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Error marshaling health check as JSON."))
		log.Println("Error marshaling health check as JSON:", err.Error())
		return
	}

	rpsls.WriteAsJson(response, healthCheckJson)
}
