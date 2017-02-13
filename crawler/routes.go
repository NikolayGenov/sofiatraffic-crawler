package crawler

type Route struct {
	Direction
	Stops
}
type Routes []Route

type Direction struct {
	Name string
	ID   string
}
type Directions []Direction

type Stop struct {
	Name        string
	CapitalName string
	Sign        string
	ID          string
	URL         string
}
type Stops []Stop
