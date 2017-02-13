package crawler

import (
	"fmt"
	"strconv"
)

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
