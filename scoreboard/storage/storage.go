package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/maxrussell/rpsls"
)

var db *sql.DB

func init() {
	// Secret management is a huge field with many diverse strategies. The approach my team and I
	// decided on at DealerPeak was:
	// * use HashiCorp's Vault to actually store the secrets
	// * each app reads, on stdin, its Vault key, which it then uses to access Vault
	// * Vault has config about which app needs which secrets and denies access to unneeded ones
	// * Vault interaction was abstracted behind one of our packages, so imported libraries
	//   practically had no way to access our secrets (and thus no way to expose them)
	//
	// Implementing this solution here seemed out of scope and also would make deploying this app
	// much more challenging, so I've decided on a simpler and popular, though less secure approach:
	// environment variables. A deployment script would just set these as it launched the server.

	host := os.Getenv("DBHOST")
	user := os.Getenv("DBUSER")
	pass := os.Getenv("DBPASS")
	database := os.Getenv("DBDATABASE")

	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s", user, database, pass, host))
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

	return tx.Commit()
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
