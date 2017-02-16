package crawler

import (
	"bytes"
	"fmt"
)

//Direction represent a line direction.
// Typically a line has 2 directions, but there could be more e.g 3 - 8.
type Direction struct {
	//Name is the name of the route - typically the names of the first
	// and last stops separated by a dash e.g "Ж.к.Западен парк - Метростанция Витоша".
	// But there are really strange examples including middle stops.
	Name string `json:"name"`
	//Unique ID for the direction of the line disregarding operation.
	// The ID is taken from schedules.sofiatraffic.bg.
	ID string `json:"id"`
}

//Stop represents the main data for a traffic stop with two names
// one capital and one normal extracted from two sources, a sign, ID
// and a entry for the same stop as a virtual tables stop if one exists.
type Stop struct {
	//Regular name of a of a stop e.g `Метростанция "Витоша"`.
	// Extracted from m.sofiatraffic.bg/schedules/.
	Name string `json:"name"`

	//Capital name of the same stop e.g  `МЕТРОСТАНЦИЯ ВИТОША`.
	// Extracted from schedules.sofiatraffic.bg.
	CapitalName string `json:"capital_name"`

	//Sign represents an ID which is marked on each actual stop sign in the real life.
	// It is second most common way to refer to a stop after its name, e.g 0910.
	// This is bridging matcher between  schedules.sofiatraffic.bg and m.sofiatraffic.bg/schedules/.
	Sign string `json:"sign"`

	//ID is a unique ID for a stop.
	// Extracted from schedules.sofiatraffic.bg.
	ID string `json:"id"`

	//VirtualTableStop is a mapped (if that mapping exist) virtual tables entry.
	// It can be used to query the given stop for real time data such as comma separated times.
	VirtualTableStop `json:"vt_stop"`
}

//Stops is a simple list of stops.
// It was created as alias to simplify printing and usage instead of slice.
type Stops []Stop

//Route is a concept which is composed of a direction and list of stops.
// Name of the route is Name of the direction.
type Route struct {
	Direction `json:"direction"`
	Stops     `json:"stops"`
}

//Routes is simple list of routes.
// It was created as alias to simplify printing and usage instead of slice.
type Routes []Route

func (s *Stop) String() string {
	return fmt.Sprintf("[%v] %v (%v)", s.Sign, s.CapitalName, s.ID)
}

func (s Stops) String() string {
	var buffer bytes.Buffer
	for _, stop := range s {
		buffer.WriteString(fmt.Sprintln(&stop))
	}
	return buffer.String()
}

func (r Route) String() string {
	return fmt.Sprintf("\n%v (%v)\n%v", r.Name, r.ID, r.Stops)
}

func (r Routes) String() string {
	var buffer bytes.Buffer
	for _, route := range r {
		buffer.WriteString(fmt.Sprint(route))
	}
	return buffer.String()
}
