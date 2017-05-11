package storage

import (
	"errors"
	"os"
	"time"

	"github.com/maxrussell/rpsls"
)

var startTime time.Time
var hostName string

func init() {
	var err error
	hostName, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	startTime = time.Now()
}

func AddResult(winner, loser rpsls.Player) error {
	return errors.New("Not yet implemented")
}

func GetTopPlayers(count int) ([]rpsls.Player, error) {
	return []rpsls.Player{}, errors.New("Not yet implemented")
}

func GetHealthCheck() rpsls.HealthCheck {
	return rpsls.HealthCheck{
		Status:    "Not Implemented",
		Name:      hostName,
		Version:   "1.0",
		StartTime: startTime.String(),             // TODO: convert to correct format
		UpTime:    time.Since(startTime).String(), // TODO: convert to correct format
	}
}
