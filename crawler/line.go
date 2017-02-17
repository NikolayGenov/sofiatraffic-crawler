package crawler

import (
	"bytes"
	"fmt"
	"strconv"
)

//OperationIDRoutesMap maps Line OperationID to list of Line Routes
type OperationIDRoutesMap map[OperationID]Routes

//Line contains all the useful (and not so useful) information about a public transportation line
type Line struct {
	//Name is the name of the line - e.g "85", "44-Б", "7-А", etc.
	Name string `json:"name"`

	//Transportation is denoting the type of Transportation of a line e.g Tram.
	Transportation `json:"transportation_type"`

	//OperationIDMap is mapping between Operation and OperationID.
	// This is needed because each line has different number of Operation modes.
	OperationIDMap `json:"operation_id_map"`

	//OperationIDRoutesMap is entry point to rest of the data for a given line
	// namely a list of all of its routes.
	OperationIDRoutesMap `json:"operation_routes_map"`
}

//scheduleIDs creates a list of all possible valid ScheduleIDs for a line by iterating trough
// all the line internal data
func (l *Line) scheduleIDs() []ScheduleID {
	scheduleIDs := make([]ScheduleID, 0)
	for operationID, routes := range l.OperationIDRoutesMap {
		for _, route := range routes {
			for _, stop := range route.Stops {
				stopID, _ := strconv.Atoi(stop.Sign)
				scheduleID := ScheduleID(fmt.Sprintf("%v/%v/%v", operationID, route.Direction.ID, stopID))
				scheduleIDs = append(scheduleIDs, scheduleID)
			}
		}
	}
	return scheduleIDs
}

func (o *OperationIDRoutesMap) String() string {
	var buffer bytes.Buffer
	for id, routes := range *o {
		buffer.WriteString(fmt.Sprintf("(OperationID: %v)%v\n", id, routes))
	}
	return buffer.String()
}

func (l *Line) String() string {
	return fmt.Sprintf("%v '%v'\n%v%v", l.Transportation, l.Name, &l.OperationIDMap, &l.OperationIDRoutesMap)
}
