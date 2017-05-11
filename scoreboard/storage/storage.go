package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/maxrussell/rpsls"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=maxrussell dbname=rpsls password=maxrussell host=localhost")
	if err != nil {
		panic(err)
	}
}

func AddResult(winner, loser rpsls.Player) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = ensureExists(winner.UserName, tx)
	if err != nil {
		return err
	}
	err = ensureExists(loser.UserName, tx)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE scoreboard SET score=score+1 WHERE username=$1", winner.UserName)
	if err != nil {
		return fmt.Errorf("Error incrementing user's score: %s", err.Error())
	}

	err = tx.Commit()
	return err
}

func GetTopPlayers(count int) ([]rpsls.Player, error) {
	rows, err := db.Query("SELECT username, score FROM scoreboard ORDER BY score DESC LIMIT $1", count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []rpsls.Player{}
	for rows.Next() {
		var name string
		var score int
		if err := rows.Scan(&name, &score); err != nil {
			return nil, err
		}

		results = append(results, rpsls.Player{
			UserName: name,
			Score:    score,
		})
	}

	return results, nil
}

func GetHealthCheck() rpsls.HealthCheck {
	var status string
	var message string
	if err := db.Ping(); err != nil {
		status = "NOT OK"
		message = err.Error()
	} else {
		status = "OK"
		message = "Connected to database postgres://localhost:5433"
	}

	return rpsls.HealthCheck{
		Status:    status,
		Name:      rpsls.HostName(),
		Version:   "1.0",
		StartTime: rpsls.StartTimeString(),
		UpTime:    rpsls.UpTimeString(),
		Dependencies: []rpsls.Dependency{rpsls.Dependency{
			Name:    "postgres",
			Status:  status,
			Message: message,
		}},
	}
}

func ensureExists(userName string, tx *sql.Tx) error {
	_, err := tx.Exec(`INSERT INTO scoreboard ( username, score ) SELECT $1, 0 WHERE NOT EXISTS( SELECT NULL FROM scoreboard WHERE username=$2 )`, userName, userName)
	if err != nil {
		return fmt.Errorf("Error ensuring user exists: %s", err.Error())
	}

	return nil
}
