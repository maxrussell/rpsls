package rpsls

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
