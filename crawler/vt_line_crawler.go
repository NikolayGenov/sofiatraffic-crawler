package crawler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/NikolayGenov/gocrawl" //Because this one ignores robots.txt
	"github.com/PuerkitoBio/goquery"
)

const (
	//schedulesURLRegexTemplate is used to match Transportation as int and Line.Name from a URL
	schedulesURLRegexTemplate = `\?tt=(.+)&ln=(.+)&s=`

	//stopRegexTemplate is used to match Stop.Name and Stop.Sign from Virtual Tables line page
	stopRegexTemplate = `(.*) \((.+)\)`

	//directionTrimmingRegexTemplate is used to remove all prefixes, suffixes that can be remove without using information
	// about the direction name. This is done because of the really really different data from schedules.sofiatraffic.bg
	// and m.sofiatraffic.bg. This does NOT solve all the cases - there are plenty of typos that can NOT be fixed.
	directionTrimmingRegexTemplate = `(Ж[.]?[ ]?К\.|ПЛ(\.|ОЩАД[А]?)|С(\.|ЕЛО)|ГР(\.|АД)|-УХО|ЦЕНТЪРА|КРАЯ НА |БУЛ. |ЛИФТ.|УЛ.|КВ[. ]?|[\- ."]+)`
)

var (
	//vtSchedulesRegex is compiled regex from schedulesURLRegexTemplate
	vtSchedulesRegex = regexp.MustCompile(schedulesURLRegexTemplate)

	//vtStopRegex is compiled regex from stopRegexTemplate
	vtStopRegex = regexp.MustCompile(stopRegexTemplate)

	//directionTrimmingRegex is compiled regex from directionTrimmingRegexTemplate
	directionTrimmingRegex = regexp.MustCompile(directionTrimmingRegexTemplate)
)

//VirtualTableStop is used to make a request for real-time-sh times for a given stop
// It is saved in Stop.VirtualTableStop, and there is a list of those only for query purposes
type VirtualTableStop struct {
	//Unique ID for a stop on m.sofiatraffic.bg logic
	StopID string `json:"stop"`

	//Line ID for the given stop on m.sofiatraffic.bg logic
	LineID string `json:"lid"`

	//Route ID for the given line on m.sofiatraffic.bg logic
	RouteID string `json:"rid"`

	//TransportationType matches Line.Transportation but is kept here for easy string access
	TransportationType string `json:"vt"`
}

func (v VirtualTableStop) String() string {
	return fmt.Sprintf("%v/%v/%v/%v", v.TransportationType, v.RouteID, v.LineID, v.StopID)
}

//vtLineCrawler is extension to gocrawl. It takes a slice of Lines, and a reference to a
// slice of Virtual Table stops, and  which is going to be filled by the crawler, and operation
// It uses operation because m.sofiatraffic.bg does not have all the data, and we know better
// that some of the lines do not operate on certain operation modes.
// It also uses a mutex to prevent potential race condition if the crawler is to be concurrent
type vtLineCrawler struct {
	gocrawl.DefaultExtender
	Operation
	Lines   []Line
	vtLines *[]VirtualTableStop
	mutex   *sync.Mutex
}

//newLinesCrawler takes a slice of lines, a creates a reference to a
// slice of Virtual Table stops and operation type and returns an initialized gocrawl.Crawler
// which is not the original gocrawl but a impolite crawler which ignores robots.txt
// It also sets  proper user agent, delay and log options
func newVirtualTableLineCrawler(lines []Line, crawlLInes *[]VirtualTableStop, operation Operation) *gocrawl.Crawler {
	vtCr := &vtLineCrawler{
		Lines:     lines,
		vtLines:   crawlLInes,
		Operation: operation,
		mutex:     &sync.Mutex{}}
	opts := gocrawl.NewOptions(vtCr)
	opts.UserAgent = userAgent
	opts.CrawlDelay = 0
	opts.MaxVisits = 0
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = true
	return gocrawl.NewCrawlerWithOptions(opts)
}

//Visit for Virtual tables schedules have the following pipeline
// - Finds which line its crawling by parsing URL to Line.Name and Line.Transportation
// - Finds a line by matching the line name from all known lines
// - Finds which is the OperationID of the line by setting the operation mode which was given to the crawler
// - Checks if the line is available on the operation mode which was requested
// - Checks if the page has returned an error which can not be explained by the information available
// - For each route:
// 	- Finds a matching route by checking a specially trimmed version of the direction name
// 	- Extract LineID, RouteID and again transportation type but this time as string
// 	- For each stop:
// 		- Finds a matching stop by comparing Stop.Sign-s
//		- Extract last StopID and creates VirtualTableStop structure
//		- Updates Stop.Name - now it is available for the first time
//		- Add VirtualTableStop to a local slice
// - Add the local slice to a global crawl slice
func (v *vtLineCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	url := ctx.URL().String()
	transportation, lineName, err := parseURLToLineParams(url)
	if err != nil {
		log.Printf("There was an error parsing [%v] to line type and name, error: %v\n", url, err)
		return nil, false
	}
	line, err := findLine(v.Lines, transportation, lineName)
	if err != nil {
		log.Printf("There was an error finding line [%v %v], error: %v\n", transportation, lineName, err)
		return nil, false
	}
	operationID, ok := line.OperationIDMap[v.Operation]
	if !ok {
		log.Printf("Line [%v %v] is not available during [%v]", transportation, lineName, v.Operation)
		return nil, false
	}
	errorMessage := strings.TrimSpace(doc.Find(".error").Text())
	if errorMessage != "" {
		log.Printf("Sofia Traffic Virtual Tables gave this error message: %v", errorMessage)
		return nil, false
	}

	routes := line.OperationIDRoutesMap[operationID]
	vtRoutesStops := make([]VirtualTableStop, 0)

	doc.Find("form").Each(func(i int, routeSelection *goquery.Selection) {
		directionString := strings.TrimSpace(routeSelection.Find(".info").Text())
		route, err := findRoute(routes, directionString)
		if err != nil {
			log.Printf("There was an error finding route with name [%v], error: %v\n", directionString, err)
			return
		}
		var lineID, routeID, transportationType string
		if lineID, err = getValueFromInput(routeSelection, "lid"); err != nil {
			log.Printf("Error input lineID on route [%v] with error: %v", route.Name, err)
			return
		}
		if routeID, err = getValueFromInput(routeSelection, "rid"); err != nil {
			log.Printf("Error input routeID on route [%v] with error: %v", route.Name, err)
			return
		}
		if transportationType, err = getValueFromInput(routeSelection, "vt"); err != nil {
			log.Printf("Error input vt on route [%v] with error: %v", route.Name, err)
			return
		}
		routeSelection.Find("option").Each(func(i int, stopSelection *goquery.Selection) {
			stopText := strings.TrimSpace(stopSelection.Text())
			stopMatches := vtStopRegex.FindStringSubmatch(stopText)
			if len(stopMatches) != 3 {
				log.Printf("Stop name is NOT in the required format: `Some name (xxxx)`, given: %v", stopText)
				return
			}
			stopName := stopMatches[1]
			stopSign := stopMatches[2]

			vtStopID, ok := stopSelection.Attr("value")
			if !ok {
				log.Printf("No required attribute 'value' found on stopID [%v (%v)]", stopName, stopSign)
			}

			stop, err := findStop(route.Stops, stopSign)
			if err != nil {
				log.Printf("Error finding stopID with name [%v (%v)] on direction [%v], error: %v\n", stopName, stopSign, route.Name, err)
				return
			}

			vtStop := VirtualTableStop{
				StopID:             vtStopID,
				LineID:             lineID,
				TransportationType: transportationType,
				RouteID:            routeID}

			stop.VirtualTableStop = vtStop
			//Update name because it is the empty string for now
			stop.Name = stopName
			vtRoutesStops = append(vtRoutesStops, vtStop)
		})
	})

	v.mutex.Lock()
	*v.vtLines = append(*v.vtLines, vtRoutesStops...)
	v.mutex.Unlock()
	return nil, false
}
func parseURLToLineParams(url string) (Transportation, string, error) {
	match := vtSchedulesRegex.FindStringSubmatch(url)
	transportationString := match[1]
	lineName := match[2]
	transportationInt, err := strconv.Atoi(transportationString)
	transportation := Transportation(transportationInt)

	return transportation, lineName, err
}

func findLine(lines []Line, transportation Transportation, name string) (*Line, error) {
	//Because in Virtual Tables there isn't a line 21-22 like in schedules.sofiatraffic.g
	if name == "21-22" {
		name = "22"
	}
	for i, line := range lines {
		if line.Transportation == transportation && line.Name == name {
			return &lines[i], nil
		}
	}
	return nil, fmt.Errorf("Line of type [%v] and name [%v] is NOT matching any known lines", transportation, name)
}

func findRoute(routes Routes, directionName string) (*Route, error) {
	trimmedDirectionName := trimDirectionName(directionName)
	for i, route := range routes {
		if trimDirectionName(route.Name) == trimmedDirectionName {
			return &routes[i], nil
		}
	}
	return nil, fmt.Errorf("The name [%v] is NOT matching any of the known direction names for this line", directionName)
}

func trimDirectionName(directionName string) string {
	upper := strings.ToUpper(directionName)
	return directionTrimmingRegex.ReplaceAllString(upper, "")
}

func getValueFromInput(s *goquery.Selection, name string) (string, error) {
	matcher := fmt.Sprintf(`input[name="%v"]`, name)
	inputSelection := s.Find(matcher).Get(0)
	for _, attr := range inputSelection.Attr {
		if attr.Key == "value" {
			return attr.Val, nil
		}
	}
	return "", fmt.Errorf("Can NOT find input tag with [%v] name", name)
}

func findStop(stops Stops, sign string) (*Stop, error) {
	for i, stop := range stops {
		if stop.Sign == sign {
			return &stops[i], nil
		}
	}
	return nil, fmt.Errorf("The sign [%v] is NOT matching any of the known stopID signs for this route", sign)
}
