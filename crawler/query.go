package crawler

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
)

func getOperationsMap(doc *goquery.Document) OperationIDMap {
	m := make(OperationIDMap)
	doc.Find(".schedule_active_list_tabs li a").
		Each(func(i int, operationSelection *goquery.Selection) {
			operationString := strings.TrimSpace(operationSelection.Text())
			operation, err := convertToOperation(operationString)
			if err != nil {
				panic(fmt.Errorf("Line MUST have of required operation types in order to be processed, given: %v", operationString))
			}
			m[operation] = getOperationIDFromElementID(operationSelection)
		})
	return m
}

func getOperationIDFromElementID(operationSelection *goquery.Selection) OperationID {
	id, ok := operationSelection.Attr("id")
	if !ok {
		html, _ := operationSelection.Html()
		panic(fmt.Errorf("Operation element MUST have 'id' element, HTML: %v", html))
	}
	parts := strings.Split(id, "_")
	if len(parts) < 2 {
		panic(fmt.Errorf("The given id '%v' is not in the required format : 'schedule_{operationID}_{somethingElse}'", id))
	}
	return OperationID(parts[1])
}

func getOperationIDRoutesMap(doc *goquery.Document) OperationIDRoutesMap {
	m := make(OperationIDRoutesMap)
	doc.Find(".schedule_active_list_content").
		Each(func(i int, operationSelection *goquery.Selection) {
			operationID := getOperationIDFromElementID(operationSelection)
			m[operationID] = getRoutes(operationSelection)
		})
	return m
}

func getRoutes(operationSelection *goquery.Selection) Routes {
	routes := Routes{}
	directorySelection := operationSelection.Find(".schedule_view_direction_tabs")
	directions := getDirections(directorySelection)
	for _, direction := range directions {
		stops := getDirectionStops(operationSelection, direction.ID)
		routes = append(routes, Route{direction, stops})
	}
	return routes
}
func getDirections(directionsSelection *goquery.Selection) Directions {
	directions := Directions{}
	directionsSelection.
		Find("a").
		Each(func(i int, directionLink *goquery.Selection) {
			name, url := getTextAndHref(directionLink)
			parts := strings.Split(url, "/")
			id := parts[len(parts)-1]
			directions = append(directions, Direction{name, id})
		})
	return directions
}

func getDirectionStops(operationSelection *goquery.Selection, directionID string) Stops {
	stops := Stops{}
	stopsMatcherString := fmt.Sprintf("[id$=_direction_%v_signs] li", directionID)
	stopsMatcher := cascadia.MustCompile(stopsMatcherString)
	operationSelection.
		FindMatcher(stopsMatcher).
		Each(func(i int, stopSelection *goquery.Selection) {
			name, _ := getTextAndHref(stopSelection.Find(".stop_change"))
			sign, url := getTextAndHref(stopSelection.Find(".stop_link"))
			url, id := getStopURLAndID(url)
			stops = append(stops, Stop{"", name, sign, id, url})
		})
	return stops
}

func getStopURLAndID(url string) (string, string) {
	//URL should be in this form : 'stop/{stop_id}/{stop_latin_name}#{integer_sign_name}'
	firstSlash := strings.Index(url, "/")
	lastSlash := strings.LastIndex(url, "/")
	url = url[:lastSlash]
	id := url[firstSlash+1:]
	return url, id
}

func getScheduleTimes(doc *goquery.Document) ScheduleTimes {
	scheduleTimes := ScheduleTimes{}
	doc.Find(".schedule_times tbody a").
		Each(func(i int, s *goquery.Selection) {
			scheduleTimes = append(scheduleTimes, strings.TrimSpace(s.Text()))
		})
	return scheduleTimes
}
