package crawler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/NikolayGenov/gocrawl" //Because this one ignores robots.txt
	"github.com/PuerkitoBio/goquery"
)

const (
	vt_schedules_url_regex      = `\?tt=(.+)&ln=(.+)&s=`
	vt_stop_regex               = `(.*) \((.+)\)`
	vt_direction_trimming_regex = `(Ж[.]?[ ]?К\.|ПЛ(\.|ОЩАД[А]?)|С(\.|ЕЛО)|ГР(\.|АД)|-УХО|ЦЕНТЪРА|КРАЯ НА |БУЛ. |ЛИФТ.|УЛ.|КВ[. ]?|[\- ."]+)`
)

var (
	vtSchedulesRegex       = regexp.MustCompile(vt_schedules_url_regex)
	vtStopRegex            = regexp.MustCompile(vt_stop_regex)
	directionTrimmingRegex = regexp.MustCompile(vt_direction_trimming_regex)
)

type VirtualTableStop struct {
	StopID             string         `json:"stop"`
	LineID             string         `json:"lid"`
	RouteID            string         `json:"rid"`
	TransportationType Transportation `json:"vt"`
}

type vtLineCrawler struct {
	gocrawl.DefaultExtender
	Operation
	Lines   *[]Line
	vtLines *[]VirtualTableStop
}

func findLine(lines []Line, transportation Transportation, name string) (*Line, error) {
	//Because in VT there is not line 21-22
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

func findStop(stops Stops, sign string) (*Stop, error) {
	for i, stop := range stops {
		if stop.Sign == sign {
			return &stops[i], nil
		}
	}
	return nil, fmt.Errorf("The sign [%v] is NOT matching any of the known stopID signs for this route", sign)
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

func trimDirectionName(directionName string) string {
	upper := strings.ToUpper(directionName)
	return directionTrimmingRegex.ReplaceAllString(upper, "")
}

func parseURLToLineParams(url string) (Transportation, string, error) {
	match := vtSchedulesRegex.FindStringSubmatch(url)
	transportationString := match[1]
	lineName := match[2]
	transportationInt, err := strconv.Atoi(transportationString)
	transportation := Transportation(transportationInt)

	return transportation, lineName, err
}
func (v *vtLineCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	url := ctx.URL().String()
	transportation, lineName, err := parseURLToLineParams(url)
	if err != nil {
		log.Printf("There was an error parsing [%v] to line type and name, error: %v\n", url, err)
		return nil, false
	}
	line, err := findLine(*v.Lines, transportation, lineName)
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
		var lineID, routeID string
		if lineID, err = getValueFromInput(routeSelection, "lid"); err != nil {
			log.Printf("Error input lineID on route [%v] with error: %v", route.Name, err)
			return
		}
		if routeID, err = getValueFromInput(routeSelection, "rid"); err != nil {
			log.Printf("Error input routeID on route [%v] with error: %v", route.Name, err)
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
				TransportationType: transportation,
				RouteID:            routeID}

			stop.VirtualTableStop = vtStop
			//Update name because it is the empty string for now
			stop.Name = stopName

			vtRoutesStops = append(vtRoutesStops, vtStop)
		})
	})

	//TODO - potential RACE CONDITION
	*v.vtLines = append(*v.vtLines, vtRoutesStops...)
	return nil, false
}

func newVirtualTableLineCrawler(lines *[]Line, crawlLInes *[]VirtualTableStop, operation Operation) crawlable {
	vtCr := &vtLineCrawler{Lines: lines, vtLines: crawlLInes, Operation: operation}
	opts := gocrawl.NewOptions(vtCr)
	opts.UserAgent = user_agent
	opts.CrawlDelay = 0
	opts.MaxVisits = 0
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = true
	c := gocrawl.NewCrawlerWithOptions(opts)
	return c
}
