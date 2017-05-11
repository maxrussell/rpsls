package rpsls

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"
)

type Player struct {
	UserName string
	Score    int
}

type HealthCheck struct {
	Status       string
	Name         string
	Version      string
	StartTime    string // consider making a type based on time.Time and using json.Marshaler
	UpTime       string // consider making a type based on time.Duration and using json.Marshaler
	Dependencies []Dependency
}

type Dependency struct {
	Name    string
	Status  string
	Message string
}

func StartTimeString() string {
	return startTime.Format(timeFormat)
}

func UpTimeString() string {
	return formatDuration(time.Since(startTime))
}

func HostName() string {
	return hostName
}

const timeFormat = "Mon January 15:04:05 2006"

func formatDuration(duration time.Duration) string {
	return fmt.Sprintf("%d:%02d:%2f", int(duration.Hours()), int(duration.Minutes())%60, math.Mod(duration.Seconds(), 60))
}

var hostName string
var startTime time.Time

func init() {
	var err error
	hostName, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	startTime = time.Now()
}

func DefaultToJson(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(response, request)
	})
}
