package crawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	//"github.com/NikolayGenov/gocrawl"
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
)

type TransportationType int

type OperationType int
type OperationTypes []OperationType

type OperationID string
type OperationTypeIDMap map[OperationType]OperationID
type OperationIDRoutesMap map[OperationID]Routes

type Line struct {
	LineBasicInfo
	OperationTypeIDMap
	OperationIDRoutesMap
}

type OperationRoutes struct {
	OperationID
	Routes
}

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

const (
	Operation_Normal OperationType = iota
	Operation_Pre_Holiday
	Operation_Holiday
)
const (
	Tram TransportationType = iota
	Trolley
	Bus
)

var (
	operationsIdentifiers = map[string]OperationType{
		"делник":                         Operation_Normal,
		"предпразник":                    Operation_Pre_Holiday,
		"празник":                        Operation_Holiday,
		"предпразник / празник":          Operation_Holiday,
		"делник / предпразник / празник": Operation_Normal}

	operationStrings = [...]string{Operation_Normal: "Weekday",
		Operation_Pre_Holiday: "Pre-Holiday",
		Operation_Holiday:     "Holiday"}

	transportationTypeStrings = [...]string{Tram: "Tram",
		Trolley: "Trolley",
		Bus:     "Bus"}
)

func (tt TransportationType) String() string {
	return transportationTypeStrings[tt]
}

func (o OperationType) String() string {
	return operationStrings[o]
}

func (l *Line) getOperationsMap(doc *goquery.Document) {
	l.OperationTypeIDMap = make(OperationTypeIDMap)
	//TODO- ^ Move that to init ?
	doc.Find(".schedule_active_list_tabs li a").
		Each(func(i int, operationSelection *goquery.Selection) {
			typeString := strings.TrimSpace(operationSelection.Text())
			operationType, ok := operationsIdentifiers[typeString]
			if !ok {
				panic(fmt.Errorf("Line MUST have of required operation types in order to be processed, given: %v", typeString))
			}
			l.OperationTypeIDMap[operationType] = getOperationIDFromElementID(operationSelection)
		})
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

func (l *Line) getOperationIDRoutesMap(doc *goquery.Document) {
	l.OperationIDRoutesMap = make(OperationIDRoutesMap)
	//TODO- ^ Move that to init ?
	doc.Find(".schedule_active_list_content").
		Each(func(i int, operationSelection *goquery.Selection) {
			operationID := getOperationIDFromElementID(operationSelection)
			l.OperationIDRoutesMap[operationID] = getRoutes(operationSelection)
		})
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

func getTextAndHref(s *goquery.Selection) (string, string) {
	text := strings.TrimSpace(s.Text())
	//TODO - maybe that is overkill ?
	if text == "" {
		html, _ := s.Html()
		panic(fmt.Errorf("The given attribute MUST have some text, HTML: %v", html))
	}

	href, ok := s.Attr("href")
	if !ok {
		html, _ := s.Html()
		panic(fmt.Errorf("The given attribute should have href attribute, HTML: %v", html))
	}
	return text, href
}

func getStopURLAndID(url string) (string, string) {
	//URL should be in this form : 'stop/{stop_id}/{stop_latin_name}#{integer_sign_name}'
	firstSlash := strings.Index(url, "/")
	lastSlash := strings.LastIndex(url, "/")
	url = url[:lastSlash]
	id := url[firstSlash+1:]
	return url, id
}

/* ========================================================================================= */
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

//Maybe move that to another file ?
type lineCrawler struct {
	gocrawl.DefaultExtender
	Line
}

func (l *Line) getLine(doc *goquery.Document) {
	l.getOperationsMap(doc)
	l.getOperationIDRoutesMap(doc)
}

func (l *lineCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	l.getLine(doc)
	return nil, false
}
func (l *LineBasicInfo) HelperTestReaderVisit(r io.Reader) Line {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	line := Line{}
	line.LineBasicInfo = *l
	line.getLine(doc)
	return line
}
