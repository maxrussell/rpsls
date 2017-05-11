package storage

import (
	"errors"

	"github.com/maxrussell/rpsls"
)

func AddResult(winner, loser rpsls.Player) error {
	return errors.New("Not yet implemented")
}

func GetTopPlayers(count int) ([]rpsls.Player, error) {
	return []rpsls.Player{}, errors.New("Not yet implemented")
}

func GetHealthCheck() rpsls.HealthCheck {
	return rpsls.HealthCheck{
		Status:    "Not Implemented",
		Name:      rpsls.HostName(),
		Version:   "1.0",
		StartTime: rpsls.StartTimeString(),
		UpTime:    rpsls.UpTimeString(),
	}
}
