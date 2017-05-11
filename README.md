# Setup
Libraries to "go get"
* github.com/lib/pq

## SQL queries to run as setup
* CREATE USER maxrussell PASSWORD 'maxrussell';
* CREATE DATABASE rpsls WITH OWNER maxrussell;
* CREATE TABLE scoreboard( username varchar(12) PRIMARY KEY NOT NULL, score int NOT NULL );
* GRANT SELECT, UPDATE, INSERT ON scoreboard TO maxrussell;

# Running The Game
* Build the two applications:
  * go install game/game.go
  * go install scoreboard/scoreboard.go
* Run the two applications:
  * $GOBIN/scoreboard
  * SCOREBOARD='http://localhost:8081';$GOBIN/game
* For the sake of time saving, the game serves HTTP on port 8080 and the scoreboard on 8081
* Check out the available choices by hitting http://localhost:8080/choices
* Play the game a bit by making requests to the game in your favorite HTTP client. Personally, I use Postman.
For example: http://localhost:8080/play?player=grace&computer=ada&choice=1
* Check out the scoreboard by hitting http://localhost:8081/scoreboard

The environment variable SCOREBOARD is optional and tells the game where the scoreboard is. If found,
the game automatically posts results to the scoreboard. This environment variable's value should be of the form
"http://localhost:8081". If not supplied, results will need to be posted to the scoreboard separately.
