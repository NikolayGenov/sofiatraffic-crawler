package crawler

import (
	"bytes"
	"fmt"
	"strconv"
)

type OperationIDMap map[Operation]OperationID

type OperationIDRoutesMap map[OperationID]Routes

type Line struct {
	LineBasicInfo
	OperationIDMap
	OperationIDRoutesMap
}

func (l *Line) LinksToCrawl(baseURL string) []string {
	links := make([]string, 0)
	for operationID, routes := range l.OperationIDRoutesMap {
		for _, route := range routes {
			for _, stop := range route.Stops {
				stopID, _ := strconv.Atoi(stop.Sign)
				links = append(links,
					fmt.Sprintf("%v/%v/%v/%v", baseURL, operationID, route.Direction.ID, stopID))
			}
		}
	}
	return links
}

func (o OperationIDMap) String() string {
	s := "["
	for operation, id := range o {
		s += fmt.Sprintf("%v (%v) ", operation, id)
	}
	s += "]"
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
	return fmt.Sprintf("%v\n%v\n%v", l.LineBasicInfo, l.OperationIDMap, l.OperationIDRoutesMap)
}
