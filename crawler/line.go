package crawler

import (
	"bytes"
	"fmt"
	"strconv"
)

type OperationIDMap map[Operation]OperationID

type OperationIDRoutesMap map[OperationID]Routes

type Line struct {
	Name string
	Path string
	Transportation
	OperationIDMap
	OperationIDRoutesMap
}

func (l *Line) ScheduleIDs() []ScheduleID {
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

func (o OperationIDMap) String() string {
	s := ""
	for operation, id := range o {
		s += fmt.Sprintf("%v (%v)\n", operation, id)
	}
	return s
}
func (o OperationIDRoutesMap) String() string {
	var buffer bytes.Buffer
	for id, routes := range o {
		buffer.WriteString(fmt.Sprintf("(OperationID: %v)%v\n", id, routes))
	}
	return buffer.String()
}

func (l Line) String() string {
	return fmt.Sprintf("%v '%v'\n%v%v", l.Transportation, l.Name, l.OperationIDMap, l.OperationIDRoutesMap)
}
