package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
)

//lineCrawler is extension to gocrawl. It takes a reference to a slice of Lines
//It uses a mutex to prevent potential race condition if the crawler is to be concurrent
type lineCrawler struct {
	gocrawl.DefaultExtender
	lines *[]Line
	mutex *sync.Mutex
}

//newLinesCrawler takes a reference to a slice of lines and
//creates a new initialized instance of gocrawl.Crawler
//with proper user agent, delay and log options set
func newLinesCrawler(lines *[]Line) *gocrawl.Crawler {
	lineCrawler := &lineCrawler{
		lines: lines,
		mutex: &sync.Mutex{}}
	opts := gocrawl.NewOptions(lineCrawler)
	opts.UserAgent = userAgent
	opts.CrawlDelay = 0
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = true
	return gocrawl.NewCrawlerWithOptions(opts)
}

//Filter here is used to exclude all the URLs that does not lead to line page
func (l *lineCrawler) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	path := ctx.URL().Path
	if path == "/" ||
		strings.HasPrefix(path, "/tramway") ||
		strings.HasPrefix(path, "/trolleybus") ||
		strings.HasPrefix(path, "/autobus") {
		return true
	}

	return false
}

///Visit is executed in two passes
// - During the first pass - using the main seed URL get all the URLs and then filter what we need from them
// - After we have filtered the elements they can pass trough here. For those that are on the second pass
//   we need to find the name and type of transport form them and then we extract operation modes and map them to their
//   operation ids and extract all the routes that the line can take - (Hint: they can be more than 2)
func (l *lineCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	path := ctx.URL().Path
	if ctx.URL().Path != "/" {
		parts := strings.Split(path[1:], "/")
		name := parts[1]
		transportation, err := convertToTransportation(parts[0])
		if err != nil || name == "" {
			panic(fmt.Errorf("Unknown transporation type or empty name, given: %v", parts))
		}
		line := Line{
			Name:                 name,
			Transportation:       transportation,
			OperationIDMap:       getOperationsMap(doc),
			OperationIDRoutesMap: getOperationIDRoutesMap(doc)}
		l.mutex.Lock()
		*l.lines = append(*l.lines, line)
		l.mutex.Unlock()
	}
	return nil, true
}

//getOperationsMap finds the necessary information to create a OperationIDMap
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
func getDirections(directionsSelection *goquery.Selection) []Direction {
	directions := make([]Direction, 0)
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
			_, id := getStopURLAndID(url)
			stops = append(stops, Stop{"", name, sign, id, VirtualTableStop{}})
		})
	return stops
}

func getStopURLAndID(url string) (string, string) {
	//URL should be in this form : 'stopID/{stop_id}/{stop_latin_name}#{integer_sign_name}'
	firstSlash := strings.Index(url, "/")
	lastSlash := strings.LastIndex(url, "/")
	url = url[:lastSlash]
	id := url[firstSlash+1:]
	return url, id
}

func getTextAndHref(s *goquery.Selection) (string, string) {
	text := strings.TrimSpace(s.Text())
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
