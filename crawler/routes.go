package crawler

import (
	"bytes"
	"fmt"
)

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
	VirtualTableStop
}
type Stops []Stop

func (s *Stop) String() string {
	return fmt.Sprintf("[%v] %v (%v)", s.Sign, s.CapitalName, s.ID)
}
func (s Stops) String() string {
	str := ""
	for _, stop := range s {
		str += fmt.Sprintf("[%v](%v) ", stop.Sign, stop.ID)
	}
	return str
}

//FIXME - Long version
//func (s Stops) String() string {
//	var buffer bytes.Buffer
//	for _, stopID := range s {
//		buffer.WriteString(fmt.Sprintln(stopID))
//	}
//	return buffer.String()
//}

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
